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

package apps

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubectl/pkg/util/podutils"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1alpha1"
	"github.com/apecloud/kubeblocks/controllers/apps/components"
	cfgcore "github.com/apecloud/kubeblocks/internal/configuration/core"
	"github.com/apecloud/kubeblocks/internal/constant"
	"github.com/apecloud/kubeblocks/internal/controller/component"
	"github.com/apecloud/kubeblocks/internal/controller/factory"
	"github.com/apecloud/kubeblocks/internal/controller/graph"
	"github.com/apecloud/kubeblocks/internal/controller/model"
	rsmcore "github.com/apecloud/kubeblocks/internal/controller/rsm"
	intctrlutil "github.com/apecloud/kubeblocks/internal/controllerutil"
	"github.com/apecloud/kubeblocks/internal/generics"
)

const (
	// componentPhaseTransition the event reason indicates that the component transits to a new phase.
	componentPhaseTransition = "ComponentPhaseTransition"

	// podContainerFailedTimeout the timeout for container of pod failures, the component phase will be set to Failed/Abnormal after this time.
	podContainerFailedTimeout = 10 * time.Second

	// podScheduledFailedTimeout timeout for scheduling failure.
	podScheduledFailedTimeout = 30 * time.Second
)

// ComponentStatusTransformer computes the current status: read the underlying rsm status and update the component status
type ComponentStatusTransformer struct {
	client.Client
}

// componentStatusHandler handles the component status
type componentStatusHandler struct {
	cli            client.Client
	reqCtx         intctrlutil.RequestCtx
	cluster        *appsv1alpha1.Cluster
	comp           *appsv1alpha1.Component
	synthesizeComp *component.SynthesizedComponent
	dag            *graph.DAG

	// runningRSM is a snapshot of the rsm that is already running
	runningRSM *workloads.ReplicatedStateMachine
	// protoRSM is the rsm object that is rebuilt from scratch during each reconcile process
	protoRSM *workloads.ReplicatedStateMachine
	// podsReady indicates if the component's underlying pods are ready
	podsReady bool
}

var _ graph.Transformer = &ComponentStatusTransformer{}

func (t *ComponentStatusTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*ComponentTransformContext)
	reqCtx := intctrlutil.RequestCtx{
		Ctx:      transCtx.Context,
		Log:      transCtx.Logger,
		Recorder: transCtx.EventRecorder,
	}
	comp := transCtx.Component
	compOrig := transCtx.ComponentOrig

	// fast return
	if model.IsObjectDeleting(compOrig) {
		return nil
	}

	cluster := transCtx.Cluster
	clusterObj := cluster.DeepCopy()
	synthesizeComp := transCtx.SynthesizeComponent

	// get underlying running rsm
	var runningRSM *workloads.ReplicatedStateMachine
	runningRSMList, err := components.ListRSMOwnedByComponent(reqCtx.Ctx, t.Client, cluster.Namespace, constant.GetComponentWellKnownLabels(cluster.Name, synthesizeComp.Name))
	if err != nil {
		return err
	}
	cnt := len(runningRSMList)
	switch {
	case cnt == 0:
		transCtx.Logger.Info(fmt.Sprintf("rsm not found, generation: %d", comp.Generation))
		runningRSM = nil
	case cnt == 1:
		runningRSM = runningRSMList[0]
	default:
		return fmt.Errorf("more than one workloads found for the component, cluster: %s, component: %s, cnt: %d", cluster.Name, synthesizeComp.Name, cnt)
	}

	// build protoRSM workload
	protoRSM, err := factory.BuildRSM(clusterObj, synthesizeComp)
	if err != nil {
		return err
	}

	switch {
	case model.IsObjectUpdating(compOrig):
		transCtx.Logger.Info(fmt.Sprintf("update component status after applying resources, generation: %d", comp.Generation))
		comp.Status.ObservedGeneration = comp.Generation
	case model.IsObjectStatusUpdating(compOrig):
		// reconcile the component status and sync the component status to cluster status
		csh := newComponentStatusHandler(reqCtx, t.Client, clusterObj, comp, synthesizeComp, runningRSM, protoRSM, dag)
		if err := csh.reconcileComponentStatus(); err != nil {
			return err
		}

		comp = csh.comp
		clusterObj = csh.cluster
	}

	graphCli, _ := transCtx.Client.(model.GraphClient)
	graphCli.Status(dag, compOrig, comp)
	graphCli.Status(dag, cluster, clusterObj)

	return nil
}

// reconcileComponentStatus reconciles phase and syncs the component status to cluster status.
func (r *componentStatusHandler) reconcileComponentStatus() error {

	// reconcile the component phase.
	if err := r.reconcileComponentPhase(); err != nil {
		return err
	}

	// sync the component status to cluster status.
	if err := r.syncComponentStatusToCluster(); err != nil {
		return err
	}

	return nil
}

// reconcileComponentPhase reconciles the component phase.
func (r *componentStatusHandler) reconcileComponentPhase() error {
	if r.runningRSM == nil {
		return nil
	}

	// check if the rsm is deleting
	isDeleting := func() bool {
		return !r.runningRSM.DeletionTimestamp.IsZero()
	}()

	// check if the rsm replicas is zero
	isZeroReplica := func() bool {
		return (r.runningRSM.Spec.Replicas == nil || *r.runningRSM.Spec.Replicas == 0) && r.synthesizeComp.Replicas == 0
	}()

	// get the component's underlying pods
	pods, err := components.ListPodOwnedByComponent(r.reqCtx.Ctx, r.cli,
		r.cluster.Namespace, constant.GetComponentWellKnownLabels(r.cluster.Name, r.synthesizeComp.Name))
	if err != nil {
		return err
	}
	hasComponentPod := func() bool {
		return len(pods) > 0
	}()

	// check if the rsm is running
	isRSMRunning, err := r.isRSMRunning()
	if err != nil {
		return err
	}

	// check if all configTemplates are synced
	isAllConfigSynced, err := r.isAllConfigSynced()
	if err != nil {
		return err
	}

	// check if the component has failed pod
	hasFailedPod, messages, err := r.hasFailedPod(pods)
	if err != nil {
		return err
	}

	// check if the component scale out failed
	isScaleOutFailed, err := r.isScaleOutFailed()
	if err != nil {
		return err
	}

	// check if the volume expansion is running
	hasRunningVolumeExpansion, hasFailedVolumeExpansion, err := r.hasVolumeExpansionRunning()
	if err != nil {
		return err
	}

	// calculate if the component has failure
	hasFailure := func() bool {
		return hasFailedPod || isScaleOutFailed || hasFailedVolumeExpansion
	}()

	// check if the component is available
	isComponentAvailable, err := r.isComponentAvailable(pods)
	if err != nil {
		return err
	}

	// check if the component is in creating phase
	isInCreatingPhase := func() bool {
		phase := r.comp.Status.Phase
		return phase == "" || phase == appsv1alpha1.CreatingClusterCompPhase
	}()

	r.podsReady = false
	switch {
	case isDeleting:
		r.setComponentStatusPhase(appsv1alpha1.DeletingClusterCompPhase, nil, "component is Deleting")
	case isZeroReplica && hasComponentPod:
		r.setComponentStatusPhase(appsv1alpha1.StoppingClusterCompPhase, nil, "component is Stopping")
		r.podsReady = true
	case isZeroReplica:
		r.setComponentStatusPhase(appsv1alpha1.StoppedClusterCompPhase, nil, "component is Stopped")
		r.podsReady = true
	case isRSMRunning && isAllConfigSynced && !hasRunningVolumeExpansion:
		r.setComponentStatusPhase(appsv1alpha1.RunningClusterCompPhase, nil, "component is Running")
		r.podsReady = true
	case !hasFailure && isInCreatingPhase:
		r.setComponentStatusPhase(appsv1alpha1.CreatingClusterCompPhase, nil, "Create a new component")
	case !hasFailure:
		r.setComponentStatusPhase(appsv1alpha1.UpdatingClusterCompPhase, nil, "component is Updating")
	case !isComponentAvailable:
		r.setComponentStatusPhase(appsv1alpha1.FailedClusterCompPhase, messages, "component is Failed")
	default:
		r.setComponentStatusPhase(appsv1alpha1.AbnormalClusterCompPhase, nil, "unknown")
	}

	// update component info to pods' annotations
	if err := components.UpdateComponentInfoToPods(r.reqCtx.Ctx, r.cli, r.cluster, r.synthesizeComp, r.dag); err != nil {
		return err
	}

	// patch the current componentSpec workload's custom labels
	if err := components.UpdateCustomLabelToPods(r.reqCtx.Ctx, r.cli, r.cluster, r.synthesizeComp, r.dag); err != nil {
		r.reqCtx.Event(r.cluster, corev1.EventTypeWarning, "Component Controller PatchWorkloadCustomLabelFailed", err.Error())
		return err
	}

	return nil
}

// syncComponentStatusToCluster syncs the component status to cluster status.
func (r *componentStatusHandler) syncComponentStatusToCluster() error {
	updatePodsReady := func(ready bool) {
		_ = r.updateClusterComponentStatus("", func(status *appsv1alpha1.ClusterComponentStatus) error {
			// if ready flag not changed, don't update the ready time
			if status.PodsReady != nil && *status.PodsReady == ready {
				return nil
			}
			status.PodsReady = &ready
			if ready {
				now := metav1.Now()
				status.PodsReadyTime = &now
			}
			return nil
		})
	}

	switch r.comp.Status.Phase {
	case appsv1alpha1.DeletingClusterCompPhase:
		r.setClusterStatusPhase(appsv1alpha1.DeletingClusterCompPhase, r.comp.Status.Message, "component is Deleting")
	case appsv1alpha1.StoppingClusterCompPhase:
		r.setClusterStatusPhase(appsv1alpha1.StoppingClusterCompPhase, r.comp.Status.Message, "component is Stopping")
	case appsv1alpha1.StoppedClusterCompPhase:
		r.setClusterStatusPhase(appsv1alpha1.StoppedClusterCompPhase, r.comp.Status.Message, "component is Stopped")
	case appsv1alpha1.RunningClusterCompPhase:
		r.setClusterStatusPhase(appsv1alpha1.RunningClusterCompPhase, r.comp.Status.Message, "component is Running")
	case appsv1alpha1.CreatingClusterCompPhase:
		r.setClusterStatusPhase(appsv1alpha1.CreatingClusterCompPhase, r.comp.Status.Message, "Create a new component")
	case appsv1alpha1.UpdatingClusterCompPhase:
		r.setClusterStatusPhase(appsv1alpha1.UpdatingClusterCompPhase, r.comp.Status.Message, "component is Updating")
	case appsv1alpha1.FailedClusterCompPhase:
		r.setClusterStatusPhase(appsv1alpha1.FailedClusterCompPhase, r.comp.Status.Message, "component is Failed")
	default:
		r.setClusterStatusPhase(appsv1alpha1.AbnormalClusterCompPhase, r.comp.Status.Message, "unknown")
	}

	updatePodsReady(r.podsReady)

	r.updateClusterMembersStatus()

	return nil
}

// isComponentAvailable tells whether the component is basically available, ether working well or in a fragile state:
// 1. at least one pod is available
// 2. with latest revision
// 3. and with leader role label set
// TODO(xingran): remove the dependency of the component's workload type.
func (r *componentStatusHandler) isComponentAvailable(pods []*corev1.Pod) (bool, error) {
	if isLatestRevision, err := components.IsComponentPodsWithLatestRevision(r.reqCtx.Ctx, r.cli, r.cluster, r.runningRSM); err != nil {
		return false, err
	} else if !isLatestRevision {
		return false, nil
	}

	shouldCheckLeader := func() bool {
		return r.synthesizeComp.WorkloadType == appsv1alpha1.Consensus || r.synthesizeComp.WorkloadType == appsv1alpha1.Replication
	}()
	hasLeaderRoleLabel := func(pod *corev1.Pod) bool {
		roleName, ok := pod.Labels[constant.RoleLabelKey]
		if !ok {
			return false
		}
		for _, replicaRole := range r.runningRSM.Spec.Roles {
			if roleName == replicaRole.Name && replicaRole.IsLeader {
				return true
			}
		}
		return false
	}
	for _, pod := range pods {
		if !podutils.IsPodAvailable(pod, 0, metav1.Time{Time: time.Now()}) {
			continue
		}
		if !shouldCheckLeader {
			continue
		}
		if _, ok := pod.Labels[constant.RoleLabelKey]; ok {
			continue
		}
		if hasLeaderRoleLabel(pod) {
			return true, nil
		}
	}
	return false, nil
}

// isRunning checks if the component underlying rsm workload is running.
func (r *componentStatusHandler) isRSMRunning() (bool, error) {
	if r.runningRSM == nil {
		return false, nil
	}
	if isLatestRevision, err := components.IsComponentPodsWithLatestRevision(r.reqCtx.Ctx, r.cli, r.cluster, r.runningRSM); err != nil {
		return false, err
	} else if !isLatestRevision {
		return false, nil
	}

	// whether rsm is ready
	return rsmcore.IsRSMReady(r.runningRSM), nil
}

// isAllConfigSynced checks if all configTemplates are synced.
func (r *componentStatusHandler) isAllConfigSynced() (bool, error) {
	var (
		cmKey client.ObjectKey
		cmObj = &corev1.ConfigMap{}
	)

	if len(r.synthesizeComp.ConfigTemplates) == 0 {
		return true, nil
	}

	configurationKey := client.ObjectKey{
		Namespace: r.cluster.Namespace,
		Name:      cfgcore.GenerateComponentConfigurationName(r.cluster.Name, r.synthesizeComp.Name),
	}
	configuration := &appsv1alpha1.Configuration{}
	if err := r.cli.Get(r.reqCtx.Ctx, configurationKey, configuration); err != nil {
		return false, err
	}
	for _, configSpec := range r.synthesizeComp.ConfigTemplates {
		item := configuration.Spec.GetConfigurationItem(configSpec.Name)
		status := configuration.Status.GetItemStatus(configSpec.Name)
		// for creating phase
		if item == nil || status == nil {
			return false, nil
		}
		cmKey = client.ObjectKey{
			Namespace: r.cluster.Namespace,
			Name:      cfgcore.GetComponentCfgName(r.cluster.Name, r.synthesizeComp.Name, configSpec.Name),
		}
		if err := r.cli.Get(r.reqCtx.Ctx, cmKey, cmObj); err != nil {
			return false, err
		}
		if intctrlutil.GetConfigSpecReconcilePhase(cmObj, *item, status) != appsv1alpha1.CFinishedPhase {
			return false, nil
		}
	}
	return true, nil
}

// isScaleOutFailed checks if the component scale out failed.
func (r *componentStatusHandler) isScaleOutFailed() (bool, error) {
	if r.runningRSM == nil {
		return false, nil
	}
	if r.runningRSM.Spec.Replicas == nil {
		return false, nil
	}
	if r.synthesizeComp.Replicas <= *r.runningRSM.Spec.Replicas {
		return false, nil
	}

	// stsObj is the underlying rsm workload which is already running in the component.
	stsObj := components.ConvertRSMToSTS(r.runningRSM)
	stsProto := components.ConvertRSMToSTS(r.protoRSM)
	backupKey := types.NamespacedName{
		Namespace: stsObj.Namespace,
		Name:      stsObj.Name + "-scaling",
	}
	d, err := components.NewDataClone(r.reqCtx, r.cli, r.cluster, r.synthesizeComp, stsObj, stsProto, backupKey)
	if err != nil {
		return false, err
	}
	if status, err := d.CheckBackupStatus(); err != nil {
		return false, err
	} else if status == components.BackupStatusFailed {
		return true, nil
	}
	for i := *r.runningRSM.Spec.Replicas; i < r.synthesizeComp.Replicas; i++ {
		if status, err := d.CheckRestoreStatus(i); err != nil {
			return false, err
		} else if status == components.BackupStatusFailed {
			return true, nil
		}
	}
	return false, nil
}

// hasVolumeExpansionRunning checks if the volume expansion is running.
func (r *componentStatusHandler) hasVolumeExpansionRunning() (bool, bool, error) {
	var (
		running bool
		failed  bool
	)
	for _, vct := range r.runningRSM.Spec.VolumeClaimTemplates {
		volumes, err := r.getRunningVolumes(r.reqCtx, r.cli, vct.Name, r.runningRSM)
		if err != nil {
			return false, false, err
		}
		for _, v := range volumes {
			if v.Status.Capacity == nil || v.Status.Capacity.Storage().Cmp(v.Spec.Resources.Requests[corev1.ResourceStorage]) >= 0 {
				continue
			}
			running = true
			// TODO: how to check the expansion failed?
		}
	}
	return running, failed, nil
}

// getRunningVolumes gets the running volumes of the rsm.
func (r *componentStatusHandler) getRunningVolumes(reqCtx intctrlutil.RequestCtx, cli client.Client, vctName string,
	rsmObj *workloads.ReplicatedStateMachine) ([]*corev1.PersistentVolumeClaim, error) {
	pvcs, err := components.ListObjWithLabelsInNamespace(reqCtx.Ctx, cli, generics.PersistentVolumeClaimSignature,
		r.cluster.Namespace, constant.GetComponentWellKnownLabels(r.cluster.Name, r.synthesizeComp.Name))
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	matchedPVCs := make([]*corev1.PersistentVolumeClaim, 0)
	prefix := fmt.Sprintf("%s-%s", vctName, rsmObj.Name)
	for _, pvc := range pvcs {
		if strings.HasPrefix(pvc.Name, prefix) {
			matchedPVCs = append(matchedPVCs, pvc)
		}
	}
	return matchedPVCs, nil
}

// hasFailedPod checks if the component has failed pod.
// TODO(xingran): remove the dependency of the component's workload type.
func (r *componentStatusHandler) hasFailedPod(pods []*corev1.Pod) (bool, appsv1alpha1.ComponentMessageMap, error) {
	if isLatestRevision, err := components.IsComponentPodsWithLatestRevision(r.reqCtx.Ctx, r.cli, r.cluster, r.runningRSM); err != nil {
		return false, nil, err
	} else if !isLatestRevision {
		return false, nil, nil
	}

	var messages appsv1alpha1.ComponentMessageMap
	// check pod readiness
	hasFailedPod, msg, _ := hasFailedAndTimedOutPod(pods)
	if hasFailedPod {
		messages = msg
		return true, messages, nil
	}
	// check role probe
	if r.synthesizeComp.WorkloadType != appsv1alpha1.Consensus && r.synthesizeComp.WorkloadType != appsv1alpha1.Replication {
		return false, messages, nil
	}
	hasProbeTimeout := false
	for _, pod := range pods {
		if _, ok := pod.Labels[constant.RoleLabelKey]; ok {
			continue
		}
		for _, condition := range pod.Status.Conditions {
			if condition.Type != corev1.PodReady || condition.Status != corev1.ConditionTrue {
				continue
			}
			podsReadyTime := &condition.LastTransitionTime
			if components.IsProbeTimeout(r.synthesizeComp.Probes, podsReadyTime) {
				hasProbeTimeout = true
				if messages == nil {
					messages = appsv1alpha1.ComponentMessageMap{}
				}
				messages.SetObjectMessage(pod.Kind, pod.Name, "Role probe timeout, check whether the application is available")
			}
		}
	}
	return hasProbeTimeout, messages, nil
}

// setComponentStatusPhase sets the component phase and messages conditionally.
func (r *componentStatusHandler) setComponentStatusPhase(phase appsv1alpha1.ClusterComponentPhase, statusMessage appsv1alpha1.ComponentMessageMap, phaseTransitionMsg string) {
	updateFn := func(status *appsv1alpha1.ComponentStatus) error {
		if status.Phase == phase {
			return nil
		}
		status.Phase = phase
		if status.Message == nil {
			status.Message = statusMessage
		} else {
			for k, v := range statusMessage {
				status.Message[k] = v
			}
		}
		return nil
	}
	if err := r.updateComponentStatus(phaseTransitionMsg, updateFn); err != nil {
		panic(fmt.Sprintf("unexpected error occurred while updating component status: %s", err.Error()))
	}
}

// updateComponentStatus updates the component status by @updateFn, with additional message to explain the transition occurred.
func (r *componentStatusHandler) updateComponentStatus(phaseTransitionMsg string, updateFn func(status *appsv1alpha1.ComponentStatus) error) error {
	if updateFn == nil {
		return nil
	}
	phase := r.comp.Status.Phase
	err := updateFn(&r.comp.Status)
	if err != nil {
		return err
	}
	if phase != r.comp.Status.Phase {
		if r.reqCtx.Recorder != nil && phaseTransitionMsg != "" {
			r.reqCtx.Recorder.Eventf(r.comp, corev1.EventTypeNormal, componentPhaseTransition, phaseTransitionMsg)
		}
	}
	return nil
}

// setStatusPhase sets the cluster component phase and messages conditionally.
func (r *componentStatusHandler) setClusterStatusPhase(phase appsv1alpha1.ClusterComponentPhase, statusMessage appsv1alpha1.ComponentMessageMap, phaseTransitionMsg string) {
	updateFn := func(status *appsv1alpha1.ClusterComponentStatus) error {
		if status.Phase == phase {
			return nil
		}
		status.Phase = phase
		if status.Message == nil {
			status.Message = statusMessage
		} else {
			for k, v := range statusMessage {
				status.Message[k] = v
			}
		}
		return nil
	}
	if err := r.updateClusterComponentStatus(phaseTransitionMsg, updateFn); err != nil {
		panic(fmt.Sprintf("unexpected error occurred while updating component status: %s", err.Error()))
	}
}

// updateClusterComponentStatus updates the cluster component status by @updateFn, with additional message to explain the transition occurred.
func (r *componentStatusHandler) updateClusterComponentStatus(phaseTransitionMsg string, updateFn func(status *appsv1alpha1.ClusterComponentStatus) error) error {
	if updateFn == nil {
		return nil
	}

	status := r.getClusterComponentStatus()
	phase := status.Phase
	err := updateFn(&status)
	if err != nil {
		return err
	}
	r.cluster.Status.Components[r.synthesizeComp.Name] = status

	if phase != status.Phase {
		// TODO: logging the event
		if r.reqCtx.Recorder != nil && phaseTransitionMsg != "" {
			r.reqCtx.Recorder.Eventf(r.cluster, corev1.EventTypeNormal, componentPhaseTransition, phaseTransitionMsg)
		}
	}

	return nil
}

// getComponentStatus gets the cluster component status.
func (r *componentStatusHandler) getClusterComponentStatus() appsv1alpha1.ClusterComponentStatus {
	if r.cluster.Status.Components == nil {
		r.cluster.Status.Components = make(map[string]appsv1alpha1.ClusterComponentStatus)
	}
	if _, ok := r.cluster.Status.Components[r.synthesizeComp.Name]; !ok {
		r.cluster.Status.Components[r.synthesizeComp.Name] = appsv1alpha1.ClusterComponentStatus{}
	}
	return r.cluster.Status.Components[r.synthesizeComp.Name]
}

// updateClusterMembersStatus updates the cluster members status.
// TODO(xingran): remove the dependency of the component's workload type.
func (r *componentStatusHandler) updateClusterMembersStatus() {
	// get component status
	componentStatus := r.getClusterComponentStatus()

	// for compatibilities prior KB 0.7.0
	buildConsensusSetStatus := func(membersStatus []workloads.MemberStatus) *appsv1alpha1.ConsensusSetStatus {
		consensusSetStatus := &appsv1alpha1.ConsensusSetStatus{
			Leader: appsv1alpha1.ConsensusMemberStatus{
				Name:       "",
				Pod:        constant.ComponentStatusDefaultPodName,
				AccessMode: appsv1alpha1.None,
			},
		}
		for _, memberStatus := range membersStatus {
			status := appsv1alpha1.ConsensusMemberStatus{
				Name:       memberStatus.Name,
				Pod:        memberStatus.PodName,
				AccessMode: appsv1alpha1.AccessMode(memberStatus.AccessMode),
			}
			switch {
			case memberStatus.IsLeader:
				consensusSetStatus.Leader = status
			case memberStatus.CanVote:
				consensusSetStatus.Followers = append(consensusSetStatus.Followers, status)
			default:
				consensusSetStatus.Learner = &status
			}
		}
		return consensusSetStatus
	}
	buildReplicationSetStatus := func(membersStatus []workloads.MemberStatus) *appsv1alpha1.ReplicationSetStatus {
		replicationSetStatus := &appsv1alpha1.ReplicationSetStatus{
			Primary: appsv1alpha1.ReplicationMemberStatus{
				Pod: "Unknown",
			},
		}
		for _, memberStatus := range membersStatus {
			status := appsv1alpha1.ReplicationMemberStatus{
				Pod: memberStatus.PodName,
			}
			switch {
			case memberStatus.IsLeader:
				replicationSetStatus.Primary = status
			default:
				replicationSetStatus.Secondaries = append(replicationSetStatus.Secondaries, status)
			}
		}
		return replicationSetStatus
	}

	// update members status
	switch r.synthesizeComp.WorkloadType {
	case appsv1alpha1.Consensus:
		componentStatus.ConsensusSetStatus = buildConsensusSetStatus(r.runningRSM.Status.MembersStatus)
	case appsv1alpha1.Replication:
		componentStatus.ReplicationSetStatus = buildReplicationSetStatus(r.runningRSM.Status.MembersStatus)
	}
	componentStatus.MembersStatus = slices.Clone(r.runningRSM.Status.MembersStatus)

	// set component status back
	r.cluster.Status.Components[r.synthesizeComp.Name] = componentStatus
}

// hasFailedAndTimedOutPod returns whether the pods of components are still failed after a PodFailedTimeout period.
func hasFailedAndTimedOutPod(pods []*corev1.Pod) (bool, appsv1alpha1.ComponentMessageMap, time.Duration) {
	var (
		hasTimedOutPod bool
		messages       = appsv1alpha1.ComponentMessageMap{}
		hasFailedPod   bool
		requeueAfter   time.Duration
	)
	for _, pod := range pods {
		isFailed, isTimedOut, messageStr := isPodFailedAndTimedOut(pod)
		if !isFailed {
			continue
		}
		if isTimedOut {
			hasTimedOutPod = true
			messages.SetObjectMessage(pod.Kind, pod.Name, messageStr)
		} else {
			hasFailedPod = true
		}
	}
	if hasFailedPod && !hasTimedOutPod {
		requeueAfter = podContainerFailedTimeout
	}
	return hasTimedOutPod, messages, requeueAfter
}

// isPodFailedAndTimedOut checks if the pod is failed and timed out.
func isPodFailedAndTimedOut(pod *corev1.Pod) (bool, bool, string) {
	if isFailed, isTimedOut, message := isPodScheduledFailedAndTimedOut(pod); isFailed {
		return isFailed, isTimedOut, message
	}
	initContainerFailed, message := isAnyContainerFailed(pod.Status.InitContainerStatuses)
	if initContainerFailed {
		return initContainerFailed, isContainerFailedAndTimedOut(pod, corev1.PodInitialized), message
	}
	containerFailed, message := isAnyContainerFailed(pod.Status.ContainerStatuses)
	if containerFailed {
		return containerFailed, isContainerFailedAndTimedOut(pod, corev1.ContainersReady), message
	}
	return false, false, ""
}

// isPodScheduledFailedAndTimedOut checks whether the unscheduled pod has timed out.
func isPodScheduledFailedAndTimedOut(pod *corev1.Pod) (bool, bool, string) {
	for _, cond := range pod.Status.Conditions {
		if cond.Type != corev1.PodScheduled {
			continue
		}
		if cond.Status == corev1.ConditionTrue {
			return false, false, ""
		}
		return true, time.Now().After(cond.LastTransitionTime.Add(podScheduledFailedTimeout)), cond.Message
	}
	return false, false, ""
}

// isAnyContainerFailed checks whether any container in the list is failed.
func isAnyContainerFailed(containersStatus []corev1.ContainerStatus) (bool, string) {
	for _, v := range containersStatus {
		waitingState := v.State.Waiting
		if waitingState != nil && waitingState.Message != "" {
			return true, waitingState.Message
		}
		terminatedState := v.State.Terminated
		if terminatedState != nil && terminatedState.Message != "" {
			return true, terminatedState.Message
		}
	}
	return false, ""
}

// isContainerFailedAndTimedOut checks whether the failed container has timed out.
func isContainerFailedAndTimedOut(pod *corev1.Pod, podConditionType corev1.PodConditionType) bool {
	containerReadyCondition := intctrlutil.GetPodCondition(&pod.Status, podConditionType)
	if containerReadyCondition == nil || containerReadyCondition.LastTransitionTime.IsZero() {
		return false
	}
	return time.Now().After(containerReadyCondition.LastTransitionTime.Add(podContainerFailedTimeout))
}

// newComponentStatusHandler creates a new componentStatusHandler
func newComponentStatusHandler(reqCtx intctrlutil.RequestCtx,
	cli client.Client,
	cluster *appsv1alpha1.Cluster,
	comp *appsv1alpha1.Component,
	synthesizeComp *component.SynthesizedComponent,
	runningRSM *workloads.ReplicatedStateMachine,
	protoRSM *workloads.ReplicatedStateMachine,
	dag *graph.DAG) *componentStatusHandler {
	return &componentStatusHandler{
		cli:            cli,
		reqCtx:         reqCtx,
		cluster:        cluster,
		comp:           comp,
		synthesizeComp: synthesizeComp,
		runningRSM:     runningRSM,
		protoRSM:       protoRSM,
		dag:            dag,
		podsReady:      false,
	}
}