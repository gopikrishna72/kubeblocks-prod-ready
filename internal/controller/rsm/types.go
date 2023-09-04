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

package rsm

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"

	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1alpha1"
	roclient "github.com/apecloud/kubeblocks/internal/controller/client"
	"github.com/apecloud/kubeblocks/internal/controller/graph"
)

const (
	// FeatureGateRSMCompatibilityMode whether run rsm in compatibility mode(i.e. make rsm having backward and forward compatibility with KubeBlocks version prior 0.6.0)
	// what the compatibility mode will do:
	// 1. labels compatibility
	//    copy labels of rsm to secondary resources(i.e. svc, sts, cm) in creation; list secondary resources by labels of rsm in update and deletion
	// 2. owner reference compatibility
	//    copy the owner ref of rsm to secondary resources; list secondary resources by owner ref of rsm in update and deletion
	// 3. finalizer compatibility
	//    copy the finalizer of rsm to secondary resources; remove the finalizer same as rsm in deletion
	// what if not:
	// labels, owner ref and finalizer of secondary resources will be generated by (not copied from) rsm.
	FeatureGateRSMCompatibilityMode = "RSM_COMPATIBILITY_MODE"

	workloadsManagedByLabelKey = "workloads.kubeblocks.io/managed-by"
	workloadsInstanceLabelKey  = "workloads.kubeblocks.io/instance"

	kindReplicatedStateMachine = "ReplicatedStateMachine"

	roleLabelKey          = "kubeblocks.io/role"
	rsmAccessModeLabelKey = "rsm.workloads.kubeblocks.io/access-mode"

	defaultPodName = "Unknown"

	rsmFinalizerName = "rsm.workloads.kubeblocks.io/finalizer"

	jobHandledLabel             = "rsm.workloads.kubeblocks.io/job-handled"
	jobTypeLabel                = "rsm.workloads.kubeblocks.io/job-type"
	jobScenarioLabel            = "rsm.workloads.kubeblocks.io/job-scenario"
	jobHandledTrue              = "true"
	jobHandledFalse             = "false"
	jobTypeSwitchover           = "switchover"
	jobTypeMemberJoinNotifying  = "member-join"
	jobTypeMemberLeaveNotifying = "member-leave"
	jobTypeLogSync              = "log-sync"
	jobTypePromote              = "promote"
	jobScenarioMembership       = "membership-reconfiguration"
	jobScenarioUpdate           = "pod-update"

	roleProbeName              = "role-observe"
	roleAgentVolumeName        = "role-agent"
	roleAgentInstallerName     = "role-agent-installer"
	roleAgentVolumeMountPath   = "/role-probe"
	roleAgentName              = "agent"
	shell2httpImage            = "msoap/shell2http:1.16.0"
	shell2httpBinaryPath       = "/app/shell2http"
	shell2httpServePath        = "/role"
	defaultRoleProbeAgentImage = "apecloud/kubeblocks-role-agent:latest"
	defaultRoleProbeDaemonPort = 7373
	roleProbeURIFormat         = "-addr=localhost:%s"
	defaultActionImage         = "busybox:latest"
	usernameCredentialVarName  = "KB_RSM_USERNAME"
	passwordCredentialVarName  = "KB_RSM_PASSWORD"
	servicePortVarName         = "KB_RSM_SERVICE_PORT"
	actionSvcListVarName       = "KB_RSM_ACTION_SVC_LIST"
	leaderHostVarName          = "KB_RSM_LEADER_HOST"
	targetHostVarName          = "KB_RSM_TARGET_HOST"
	roleProbeEventFieldPath    = "spec.containers{" + roleProbeName + "}"
	actionSvcPortBase          = int32(36500)
)

type rsmTransformContext struct {
	context.Context
	Client roclient.ReadonlyClient
	record.EventRecorder
	logr.Logger
	rsm     *workloads.ReplicatedStateMachine
	rsmOrig *workloads.ReplicatedStateMachine
}

func (c *rsmTransformContext) GetContext() context.Context {
	return c.Context
}

func (c *rsmTransformContext) GetClient() roclient.ReadonlyClient {
	return c.Client
}

func (c *rsmTransformContext) GetRecorder() record.EventRecorder {
	return c.EventRecorder
}

func (c *rsmTransformContext) GetLogger() logr.Logger {
	return c.Logger
}

var _ graph.TransformContext = &rsmTransformContext{}