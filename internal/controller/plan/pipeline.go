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

package plan

import (
	"encoding/json"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/internal/configuration/core"
	cfgutil "github.com/apecloud/kubeblocks/internal/configuration/util"
	"github.com/apecloud/kubeblocks/internal/constant"
	"github.com/apecloud/kubeblocks/internal/controller/builder"
	"github.com/apecloud/kubeblocks/internal/controller/component"
	intctrlutil "github.com/apecloud/kubeblocks/internal/controllerutil"
)

type ReconcileCtx struct {
	*intctrlutil.ResourceCtx

	Cluster    *appsv1alpha1.Cluster
	ClusterVer *appsv1alpha1.ClusterVersion
	Component  *component.SynthesizedComponent
	PodSpec    *corev1.PodSpec

	Object client.Object
	Cache  []client.Object
}

type pipeline struct {
	configuration *appsv1alpha1.Configuration
	renderWrapper renderWrapper

	ctx ReconcileCtx
	intctrlutil.ResourceFetcher[pipeline]
}

type updatePipeline struct {
	reconcile     bool
	renderWrapper renderWrapper

	item       appsv1alpha1.ConfigurationItemDetail
	itemStatus *appsv1alpha1.ConfigurationItemDetailStatus
	configSpec *appsv1alpha1.ComponentConfigSpec
	originalCM *corev1.ConfigMap
	newCM      *corev1.ConfigMap

	ctx ReconcileCtx
	intctrlutil.ResourceFetcher[updatePipeline]
}

func NewCreatePipeline(ctx ReconcileCtx) *pipeline {
	p := &pipeline{ctx: ctx}
	return p.Init(ctx.ResourceCtx, p)
}

func NewReconcilePipeline(ctx ReconcileCtx, item appsv1alpha1.ConfigurationItemDetail, itemStatus *appsv1alpha1.ConfigurationItemDetailStatus) *updatePipeline {
	p := &updatePipeline{
		reconcile:  true,
		item:       item,
		itemStatus: itemStatus,
		ctx:        ctx,
	}
	return p.Init(ctx.ResourceCtx, p)
}

func (p *pipeline) Prepare() *pipeline {
	buildTemplate := func() (err error) {
		ctx := p.ctx
		templateBuilder := newTemplateBuilder(p.ClusterName, p.Namespace, ctx.Cluster, ctx.ClusterVer, p.Context, p.Client)
		// Prepare built-in objects and built-in functions
		if err = templateBuilder.injectBuiltInObjectsAndFunctions(ctx.PodSpec, ctx.Component.ConfigTemplates, ctx.Component, ctx.Cache); err != nil {
			return
		}
		p.renderWrapper = newTemplateRenderWrapper(templateBuilder, ctx.Cluster, p.Context, ctx.Client)
		return
	}

	return p.Wrap(buildTemplate)
}

func (p *pipeline) RenderScriptTemplate() *pipeline {
	return p.Wrap(func() error {
		ctx := p.ctx
		return p.renderWrapper.renderScriptTemplate(ctx.Cluster, ctx.Component, ctx.Cache)
	})
}

func (p *pipeline) Configuration() *pipeline {
	buildConfiguration := func() (err error) {
		expectConfiguration := p.createConfiguration()
		configuration := appsv1alpha1.Configuration{}
		err = p.ResourceFetcher.Client.Get(p.Context, client.ObjectKeyFromObject(expectConfiguration), &configuration)
		switch {
		case err == nil:
			return p.updateConfiguration(&configuration, expectConfiguration)
		case !apierrors.IsNotFound(err):
			return p.ResourceFetcher.Client.Create(p.Context, expectConfiguration)
		default:
			return err
		}
	}
	return p.Wrap(buildConfiguration)
}

func (p *pipeline) CreateConfigTemplate() *pipeline {
	return p.Wrap(func() error {
		ctx := p.ctx
		return p.renderWrapper.renderConfigTemplate(ctx.Cluster, ctx.Component, ctx.Cache)
	})
}

func (p *pipeline) UpdateConfigurationStatus() *pipeline {
	return p.Wrap(func() error {
		if p.configuration != nil {
			return nil
		}

		patch := client.MergeFrom(p.configuration)
		updated := p.configuration.DeepCopy()
		for _, item := range p.configuration.Spec.ConfigItemDetails {
			checkAndUpdateItemStatus(updated, item)
		}
		return p.ResourceFetcher.Client.Status().Patch(p.Context, updated, patch)
	})
}

func checkAndUpdateItemStatus(updated *appsv1alpha1.Configuration, item appsv1alpha1.ConfigurationItemDetail) {
	foundStatus := func(name string) *appsv1alpha1.ConfigurationItemDetailStatus {
		for i := range updated.Status.ConfigurationItemStatus {
			status := &updated.Status.ConfigurationItemStatus[i]
			if status.Name == name {
				return status
			}
		}
		return nil
	}

	status := foundStatus(item.Name)
	if status == nil {
		updated.Status.ConfigurationItemStatus = append(updated.Status.ConfigurationItemStatus,
			appsv1alpha1.ConfigurationItemDetailStatus{
				Name:  item.Name,
				Phase: appsv1alpha1.CInitPhase,
			})
	}
}

func (p *pipeline) UpdatePodVolumes() *pipeline {
	return p.Wrap(func() error {
		return intctrlutil.CreateOrUpdatePodVolumes(p.ctx.PodSpec, p.renderWrapper.volumes)
	})
}

func (p *pipeline) BuildConfigManagerSidecar() *pipeline {
	return p.Wrap(func() error {
		return buildConfigManagerWithComponent(p.ctx.PodSpec, p.ctx.Component.ConfigTemplates, p.Context, p.Client, p.ctx.Cluster, p.ctx.Component)
	})
}

func (p *pipeline) UpdateConfigMeta() *pipeline {
	updateMeta := func() error {
		updateResourceAnnotationsWithTemplate(p.ctx.Object, p.renderWrapper.templateAnnotations)
		if err := injectTemplateEnvFrom(p.ctx.Cluster, p.ctx.Component, p.ctx.PodSpec, p.Client, p.Context, p.renderWrapper.renderedObjs); err != nil {
			return err
		}
		return createConfigObjects(p.Client, p.Context, p.renderWrapper.renderedObjs)
	}

	return p.Wrap(updateMeta)
}

func (p *pipeline) createConfiguration() *appsv1alpha1.Configuration {
	builder := builder.NewConfigurationBuilder(p.Namespace,
		core.GenerateComponentConfigurationName(p.ClusterName, p.ComponentName),
	)

	for _, template := range p.ctx.Component.ConfigTemplates {
		builder.AddConfigurationItem(template.Name)
	}
	return builder.Component(p.ComponentName).
		ClusterRef(p.ClusterName).
		GetObject()
}

func (p *pipeline) updateConfiguration(expected *appsv1alpha1.Configuration, existing *appsv1alpha1.Configuration) error {
	fromMap := func(items []appsv1alpha1.ConfigurationItemDetail) *cfgutil.Sets {
		sets := cfgutil.NewSet()
		for _, item := range items {
			sets.Add(item.Name)
		}
		return sets
	}

	oldSets := fromMap(existing.Spec.ConfigItemDetails)
	newSets := fromMap(expected.Spec.ConfigItemDetails)

	addSets := cfgutil.Difference(newSets, oldSets)
	delSets := cfgutil.Difference(oldSets, newSets)

	newConfigItems := make([]appsv1alpha1.ConfigurationItemDetail, 0)
	for _, item := range existing.Spec.ConfigItemDetails {
		if !delSets.InArray(item.Name) {
			newConfigItems = append(newConfigItems, item)
		}
	}
	for _, item := range expected.Spec.ConfigItemDetails {
		if addSets.InArray(item.Name) {
			newConfigItems = append(newConfigItems, item)
		}
	}

	patch := client.MergeFrom(existing)
	updated := existing.DeepCopy()
	updated.Spec.ConfigItemDetails = newConfigItems
	return p.Client.Patch(p.Context, updated, patch)
}

func (p *updatePipeline) isDone() bool {
	return !p.reconcile
}

func (p *updatePipeline) Prepare() *updatePipeline {
	buildTemplate := func() (err error) {
		p.reconcile = !isFinish(p.originalCM, p.item)
		if !p.isDone() {
			return
		}
		templateBuilder := newTemplateBuilder(p.ClusterName, p.Namespace, p.ctx.Cluster, p.ctx.ClusterVer, p.Context, p.Client)
		// Prepare built-in objects and built-in functions
		if err = templateBuilder.injectBuiltInObjectsAndFunctions(p.ctx.PodSpec, []appsv1alpha1.ComponentConfigSpec{*p.configSpec}, p.ctx.Component, p.ctx.Cache); err != nil {
			return
		}
		p.renderWrapper = newTemplateRenderWrapper(templateBuilder, p.ctx.Cluster, p.Context, p.Client)
		return
	}
	return p.Wrap(buildTemplate)
}

func isFinish(cm *corev1.ConfigMap, item appsv1alpha1.ConfigurationItemDetail) bool {
	if cm == nil {
		return false
	}

	lastAppliedVersion, ok := cm.Annotations[constant.ConfigAppliedVersionAnnotationKey]
	if !ok {
		return false
	}
	var target appsv1alpha1.ConfigurationItemDetail
	if err := json.Unmarshal([]byte(lastAppliedVersion), &target); err != nil {
		return false
	}

	return reflect.DeepEqual(target, item)
}

func (p *updatePipeline) ConfigMap() *updatePipeline {
	cmKey := client.ObjectKey{
		Namespace: p.Namespace,
		Name:      core.GetComponentCfgName(p.ClusterName, p.ComponentName, p.item.Name),
	}
	return p.Wrap(func() error {
		configSpec, err := p.foundConfigSpec(p.item.Name)
		if err != nil {
			return err
		}
		p.configSpec = configSpec
		p.originalCM = &corev1.ConfigMap{}
		return p.Client.Get(p.Context, cmKey, p.originalCM)
	})
}

func (p *updatePipeline) foundConfigSpec(configSpec string) (*appsv1alpha1.ComponentConfigSpec, error) {
	var i int
	templates := p.ctx.Component.ConfigTemplates
	for i = 0; i < len(templates); i++ {
		if templates[i].Name == configSpec {
			break
		}
	}
	if i >= len(templates) {
		return nil, core.MakeError("not found config spec: %s", configSpec)
	}
	return &templates[i], nil
}

// step1: doRerender
func (p *updatePipeline) RerenderTemplate() *updatePipeline {
	return p.Wrap(func() (err error) {
		if p.isDone() {
			return
		}
		if needRerender(p.originalCM, p.item) {
			p.newCM, err = p.renderWrapper.rerenderConfigTemplate(p.ctx.Cluster, p.ctx.Component, *p.configSpec, p.item)
		}
		return
	})
}

// step1: doMerge
func (p *updatePipeline) ApplyParameters() *updatePipeline {
	return p.Wrap(func() error {
		if p.isDone() {
			return nil
		}
		return p.renderWrapper.updateConfigTemplate(p.ctx.Cluster, p.ctx.Component, *p.configSpec, p.itemStatus)
	})
}

func (p *updatePipeline) UpdateConfigVersion() *updatePipeline {
	return p.Wrap(func() error {
		if p.isDone() {
			return nil
		}
		return p.renderWrapper.updateConfigTemplate(p.ctx.Cluster, p.ctx.Component, *p.configSpec, p.itemStatus)
	})
}

// step1: doSync
func (p *updatePipeline) Sync() *updatePipeline {
	return p.Wrap(func() error {
		if p.isDone() {
			return nil
		}
		// TODO merge and sync
		return nil
	})
}

func (p *updatePipeline) SyncStatus() *updatePipeline {
	return p.Wrap(func() error {
		if p.configSpec == nil {
			return nil
		}
		// TODO merge and sync
		return nil
	})
}
