/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dapr/kit/logger"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/apecloud/kubeblocks/cmd/probe/internal"
	"github.com/apecloud/kubeblocks/cmd/probe/internal/component"
	"github.com/apecloud/kubeblocks/cmd/probe/internal/dcs"
)

type WesqlManager struct {
	Manager
}

var _ component.DBManager = &WesqlManager{}

func NewWesqlManager(logger logger.Logger) (*WesqlManager, error) {
	db, err := config.GetLocalDBConn()
	if err != nil {
		return nil, errors.Wrap(err, "connect to MySQL")
	}

	defer func() {
		if err != nil {
			derr := db.Close()
			if derr != nil {
				logger.Errorf("failed to close: %v", err)
			}
		}
	}()

	currentMemberName := viper.GetString("KB_POD_NAME")
	if currentMemberName == "" {
		return nil, fmt.Errorf("KB_POD_NAME is not set")
	}

	serverID, err := component.GetIndex(currentMemberName)
	if err != nil {
		return nil, err
	}

	mgr := &WesqlManager{
		Manager: Manager{
			DBManagerBase: component.DBManagerBase{
				CurrentMemberName: currentMemberName,
				ClusterCompName:   viper.GetString("KB_CLUSTER_COMP_NAME"),
				Namespace:         viper.GetString("KB_NAMESPACE"),
				Logger:            logger,
			},
			DB:       db,
			serverID: uint(serverID) + 1,
		},
	}

	component.RegisterManager("mysql", internal.Consensus, mgr)
	return mgr, nil
}

func (mgr *WesqlManager) InitializeCluster(ctx context.Context, cluster *dcs.Cluster) error {
	return nil
}

func (mgr *WesqlManager) IsRunning() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// test if db is ready to connect or not
	err := mgr.DB.PingContext(ctx)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			// Now the error number is accessible directly
			if driverErr.Number == 1040 {
				mgr.Logger.Infof("Too many connections: %v", err)
				return true
			}
		}
		mgr.Logger.Infof("DB is not ready: %v", err)
		return false
	}

	return true
}

func (mgr *WesqlManager) IsDBStartupReady() bool {
	if mgr.DBStartupReady {
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// test if db is ready to connect or not
	err := mgr.DB.PingContext(ctx)
	if err != nil {
		mgr.Logger.Infof("DB is not ready: %v", err)
		return false
	}

	mgr.DBStartupReady = true
	mgr.Logger.Infof("DB startup ready")
	return true
}

func (mgr *WesqlManager) IsReadonly(ctx context.Context, cluster *dcs.Cluster, member *dcs.Member) (bool, error) {
	var db *sql.DB
	var err error
	if member != nil {
		addr := cluster.GetMemberAddrWithPort(*member)
		db, err = config.GetDBConnWithAddr(addr)
		if err != nil {
			mgr.Logger.Infof("Get Member conn failed: %v", err)
			return false, err
		}
		if db != nil {
			defer db.Close()
		}
	} else {
		db = mgr.DB
	}

	var readonly bool
	err = db.QueryRowContext(ctx, "select @@global.hostname, @@global.version, "+
		"@@global.read_only, @@global.binlog_format, @@global.log_bin, @@global.log_slave_updates").
		Scan(&mgr.hostname, &mgr.version, &readonly, &mgr.binlogFormat,
			&mgr.logbinEnabled, &mgr.logReplicationUpdatesEnabled)
	if err != nil {
		mgr.Logger.Infof("Get global readonly failed: %v", err)
		return false, err
	}
	return readonly, nil
}

func (mgr *WesqlManager) IsLeader(ctx context.Context, cluster *dcs.Cluster) (bool, error) {
	role, err := mgr.GetRole(ctx)

	if err == nil && strings.EqualFold(role, "leader") {
		return true, nil
	}

	return false, err
}

func (mgr *WesqlManager) IsLeaderMember(ctx context.Context, cluster *dcs.Cluster, member *dcs.Member) (bool, error) {
	readonly, err := mgr.IsReadonly(ctx, cluster, member)
	if err != nil || readonly {
		return false, err
	}

	return true, err
}

func (mgr *WesqlManager) InitiateCluster(cluster *dcs.Cluster) error {
	return nil
}

func (mgr *WesqlManager) GetMemberAddrs(cluster *dcs.Cluster) []string {
	return cluster.GetMemberAddrs()
}

func (mgr *WesqlManager) GetLeaderClient(ctx context.Context, cluster *dcs.Cluster) (*sql.DB, error) {
	leaderMember := cluster.GetLeaderMember()
	if leaderMember == nil {
		return nil, fmt.Errorf("cluster has no leader")
	}

	addr := cluster.GetMemberAddrWithPort(*leaderMember)
	return config.GetDBConnWithAddr(addr)
}

func (mgr *WesqlManager) IsCurrentMemberInCluster(ctx context.Context, cluster *dcs.Cluster) bool {
	clusterInfo := mgr.GetClusterInfo(ctx, cluster)
	return strings.Contains(clusterInfo, mgr.CurrentMemberName)
}

func (mgr *WesqlManager) IsCurrentMemberHealthy(ctx context.Context, cluster *dcs.Cluster) bool {
	mgr.DBState = nil
	member := cluster.GetMemberWithName(mgr.CurrentMemberName)
	if !mgr.IsMemberHealthy(ctx, cluster, member) {
		return false
	}

	mgr.DBState = mgr.GetDBState(ctx, cluster, member)
	if cluster.Leader != nil && cluster.Leader.Name == member.Name {
		cluster.Leader.DBState = mgr.DBState
	}
	return true
}

func (mgr *WesqlManager) IsMemberLagging(ctx context.Context, cluster *dcs.Cluster, member *dcs.Member) bool {
	return false
}

func (mgr *WesqlManager) IsMemberHealthy(ctx context.Context, cluster *dcs.Cluster, member *dcs.Member) bool {
	var db *sql.DB
	var err error
	if member != nil && member.Name != mgr.CurrentMemberName {
		db, err = mgr.GetDBConnWithMember(cluster, member)
		if err != nil {
			mgr.Logger.Infof("Get Member conn failed: %v", err)
			return false
		}
		if db != nil {
			defer db.Close()
		}
	} else {
		db = mgr.DB
	}

	if cluster.Leader != nil && cluster.Leader.Name == member.Name {
		if !mgr.WriteCheck(ctx, db) {
			return false
		}
	}
	if !mgr.ReadCheck(ctx, db) {
		return false
	}

	return true
}

func (mgr *WesqlManager) GetDBState(ctx context.Context, cluster *dcs.Cluster, member *dcs.Member) *dcs.DBState {
	var db *sql.DB
	var err error
	var isCurrentMember bool
	if member != nil && member.Name != mgr.CurrentMemberName {
		addr := cluster.GetMemberAddrWithPort(*member)
		db, err = config.GetDBConnWithAddr(addr)
		if err != nil {
			mgr.Logger.Infof("Get Member conn failed: %v", err)
			return nil
		}
		if db != nil {
			defer db.Close()
		}
	} else {
		isCurrentMember = true
		db = mgr.DB
	}

	globalState, err := mgr.GetGlobalState(ctx, db)
	if err != nil {
		mgr.Logger.Infof("select global failed: %v", err)
		return nil
	}

	opTimestamp, err := mgr.GetOpTimestamp(ctx, db)
	if err != nil {
		mgr.Logger.Infof("get op timestamp failed: %v", err)
		return nil
	}

	dbState := &dcs.DBState{
		OpTimestamp: opTimestamp,
		Extra:       map[string]string{},
	}

	for k, v := range globalState {
		dbState.Extra[k] = v
	}

	if isCurrentMember {
		mgr.globalState = globalState
		mgr.opTimestamp = opTimestamp
	}
	return dbState
}

func (mgr *WesqlManager) WriteCheck(ctx context.Context, db *sql.DB) bool {
	writeSQL := fmt.Sprintf(`BEGIN;
CREATE DATABASE IF NOT EXISTS kubeblocks;
CREATE TABLE IF NOT EXISTS kubeblocks.kb_health_check(type INT, check_ts BIGINT, PRIMARY KEY(type));
INSERT INTO kubeblocks.kb_health_check VALUES(%d, UNIX_TIMESTAMP()) ON DUPLICATE KEY UPDATE check_ts = UNIX_TIMESTAMP();
COMMIT;`, component.CheckStatusType)
	_, err := db.ExecContext(ctx, writeSQL)
	if err != nil {
		mgr.Logger.Infof("SQL %s executing failed: %v", writeSQL, err)
		return false
	}
	return true
}

func (mgr *WesqlManager) ReadCheck(ctx context.Context, db *sql.DB) bool {
	_, err := mgr.GetOpTimestamp(ctx, db)
	if err != nil {
		if err == sql.ErrNoRows {
			// no healthy check records, return true
			return true
		}
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1049 {
			// no healthy database, return true
			return true
		}
		mgr.Logger.Infof("Read check failed: %v", err)
		return false
	}

	return true
}

func (mgr *WesqlManager) GetOpTimestamp(ctx context.Context, db *sql.DB) (int64, error) {
	readSQL := fmt.Sprintf(`select check_ts from kubeblocks.kb_health_check where type=%d limit 1;`, component.CheckStatusType)
	var opTimestamp int64
	err := db.QueryRowContext(ctx, readSQL).Scan(&opTimestamp)
	return opTimestamp, err
}

func (mgr *WesqlManager) GetGlobalState(ctx context.Context, db *sql.DB) (map[string]string, error) {
	var hostname, serverUUID, gtidExecuted, gtidPurged, isReadonly, superReadonly string
	err := db.QueryRowContext(ctx, "select  @@global.hostname, @@global.server_uuid, @@global.gtid_executed, @@global.gtid_purged, @@global.read_only, @@global.super_read_only").
		Scan(&hostname, &serverUUID, &gtidExecuted, &gtidPurged, &isReadonly, &superReadonly)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"hostname":        hostname,
		"server_uuid":     serverUUID,
		"gtid_executed":   gtidExecuted,
		"gtid_purged":     gtidPurged,
		"read_only":       isReadonly,
		"super_read_only": superReadonly,
	}, nil
}

func (mgr *WesqlManager) Recover(context.Context) error {
	return nil
}

func (mgr *WesqlManager) AddCurrentMemberToCluster(cluster *dcs.Cluster) error {
	return nil
}

func (mgr *WesqlManager) DeleteMemberFromCluster(cluster *dcs.Cluster, host string) error {
	return nil
}

func (mgr *WesqlManager) IsClusterHealthy(ctx context.Context, cluster *dcs.Cluster) bool {
	db, err := mgr.GetLeaderConn(ctx, cluster)
	if err != nil {
		mgr.Logger.Infof("Get leader conn failed: %v", err)
		return false
	}
	if db == nil {
		return false
	}

	defer db.Close()
	var leaderRecord RowMap
	sql := "select * from information_schema.wesql_cluster_global;"
	err = QueryRowsMap(db, sql, func(rMap RowMap) error {
		if rMap.GetString("ROLE") == "Leader" {
			leaderRecord = rMap
		}
		return nil
	})
	if err != nil {
		mgr.Logger.Errorf("error executing %s: %v", sql, err)
		return false
	}

	if len(leaderRecord) > 0 {
		return true
	}
	return false
}

// IsClusterInitialized is a method to check if cluster is initailized or not
func (mgr *WesqlManager) IsClusterInitialized(ctx context.Context, cluster *dcs.Cluster) (bool, error) {
	clusterInfo := mgr.GetClusterInfo(ctx, nil)
	if clusterInfo != "" {
		return true, nil
	}

	return false, nil
}

func (mgr *WesqlManager) GetClusterInfo(ctx context.Context, cluster *dcs.Cluster) string {
	var db *sql.DB
	var err error
	if cluster != nil {
		db, err = mgr.GetLeaderConn(ctx, cluster)
		if err != nil {
			mgr.Logger.Infof("Get leader conn failed: %v", err)
			return ""
		}
		if db != nil {
			defer db.Close()
		}
	} else {
		db = mgr.DB

	}
	var clusterID, clusterInfo string
	err = db.QueryRowContext(ctx, "select cluster_id, cluster_info from mysql.consensus_info").
		Scan(&clusterID, &clusterInfo)
	if err != nil {
		mgr.Logger.Error("Cluster info query failed: %v", err)
	}
	return clusterInfo
}

func (mgr *WesqlManager) Promote(ctx context.Context, cluster *dcs.Cluster) error {
	isLeader, _ := mgr.IsLeader(ctx, cluster)
	if isLeader {
		return nil
	}

	db, err := mgr.GetLeaderConn(ctx, cluster)
	if err != nil {
		return errors.Wrap(err, "Get leader conn failed")
	}
	if db != nil {
		defer db.Close()
	}

	currentMember := cluster.GetMemberWithName(mgr.GetCurrentMemberName())
	addr := cluster.GetMemberAddr(*currentMember)
	resp, err := db.Exec(fmt.Sprintf("call dbms_consensus.change_leader('%s:13306');", addr))
	if err != nil {
		mgr.Logger.Errorf("promote err: %v", err)
		return err
	}

	mgr.Logger.Infof("promote success, resp:%v", resp)
	return nil
}

func (mgr *WesqlManager) Demote(context.Context) error {
	return nil
}

func (mgr *WesqlManager) Follow(ctx context.Context, cluster *dcs.Cluster) error {
	return nil
}

func (mgr *WesqlManager) GetHealthiestMember(cluster *dcs.Cluster, candidate string) *dcs.Member {
	return nil
}

func (mgr *WesqlManager) HasOtherHealthyLeader(ctx context.Context, cluster *dcs.Cluster) *dcs.Member {
	clusterLocalInfo, err := mgr.GetClusterLocalInfo(ctx)
	if err != nil || clusterLocalInfo == nil {
		mgr.Logger.Errorf("Get cluster local info failed: %v", err)
		return nil
	}

	if clusterLocalInfo.GetString("ROLE") == "Leader" {
		// I am the leader, just return nil
		return nil
	}

	leaderAddr := clusterLocalInfo.GetString("CURRENT_LEADER")
	if leaderAddr == "" {
		return nil
	}
	leaderParts := strings.Split(leaderAddr, ".")
	if len(leaderParts) > 0 {
		return cluster.GetMemberWithName(leaderParts[0])
	}

	return nil
}

// HasOtherHealthyMembers checks if there are any healthy members, excluding the leader
func (mgr *WesqlManager) HasOtherHealthyMembers(ctx context.Context, cluster *dcs.Cluster, leader string) []*dcs.Member {
	members := make([]*dcs.Member, 0)
	for _, member := range cluster.Members {
		if member.Name == leader {
			continue
		}
		if !mgr.IsMemberHealthy(ctx, cluster, &member) {
			continue
		}
		members = append(members, &member)
	}

	return members
}

func (mgr *WesqlManager) Lock(ctx context.Context, reason string) error {
	setReadOnly := `set global read_only=on;`

	_, err := mgr.DB.Exec(setReadOnly)
	if err != nil {
		mgr.Logger.Errorf("Lock err: %v", err)
		return err
	}
	mgr.IsLocked = true
	return nil
}

func (mgr *WesqlManager) Unlock(ctx context.Context) error {
	setReadOnlyOff := `set global read_only=off;`

	_, err := mgr.DB.Exec(setReadOnlyOff)
	if err != nil {
		mgr.Logger.Errorf("Unlock err: %v", err)
		return err
	}
	mgr.IsLocked = false
	return nil
}