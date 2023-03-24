/*
Copyright ApeCloud, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lifecycle

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/internal/controller/builder"
	"github.com/apecloud/kubeblocks/internal/controller/plan"
	intctrlutil "github.com/apecloud/kubeblocks/internal/controllerutil"
)

// TODO: define a custom workload to encapsulate all the resources.
type componentBuilder interface {
	buildEnv() componentBuilder
	buildHeadlessService() componentBuilder
	buildService() componentBuilder
	buildTLSCert() componentBuilder
	complete() error

	// workload related
	buildConfig(idx int32) componentBuilder
	buildWorkload(idx int32) componentBuilder
	buildVolume(idx int32) componentBuilder
	buildVolumeMount(idx int32) componentBuilder
	buildTLSVolume(idx int32) componentBuilder

	mutableWorkload(idx int32) client.Object
	mutablePodSpec(idx int32) *corev1.PodSpec
}

// single workload component
type componentBuilderBase struct {
	ReqCtx          intctrlutil.RequestCtx
	Client          client.Client
	Comp            Component
	defaultAction   *Action
	concreteBuilder componentBuilder
	Error           error
	EnvConfig       *corev1.ConfigMap
}

func (b *componentBuilderBase) buildEnv() componentBuilder {
	buildfn := func() ([]client.Object, error) {
		envCfg, err := builder.BuildEnvConfigLow(b.ReqCtx, b.Client, b.Comp.GetCluster(), b.Comp.GetSynthesizedComponent())
		b.EnvConfig = envCfg
		return []client.Object{envCfg}, err
	}
	return b.buildWrapper(buildfn)
}

func (b *componentBuilderBase) buildConfig(idx int32) componentBuilder {
	buildfn := func() ([]client.Object, error) {
		workload := b.concreteBuilder.mutableWorkload(idx)
		if workload == nil {
			return nil, fmt.Errorf("build config but workload is nil, cluster: %s, component: %s",
				b.Comp.GetClusterName(), b.Comp.GetName())
		}

		return plan.BuildCfgLow(b.Comp.GetVersion(), b.Comp.GetCluster(), b.Comp.GetSynthesizedComponent(), workload,
			b.concreteBuilder.mutablePodSpec(idx), b.ReqCtx.Ctx, b.Client)
	}
	return b.buildWrapper(buildfn)
}

func (b *componentBuilderBase) buildHeadlessService() componentBuilder {
	buildfn := func() ([]client.Object, error) {
		svc, err := builder.BuildHeadlessSvcLow(b.Comp.GetCluster(), b.Comp.GetSynthesizedComponent())
		return []client.Object{svc}, err
	}
	return b.buildWrapper(buildfn)
}

func (b *componentBuilderBase) buildService() componentBuilder {
	buildfn := func() ([]client.Object, error) {
		svcList, err := builder.BuildSvcListLow(b.Comp.GetCluster(), b.Comp.GetSynthesizedComponent())
		if err != nil {
			return nil, err
		}
		objs := make([]client.Object, len(svcList))
		for _, svc := range svcList {
			objs = append(objs, svc)
		}
		return objs, err
	}
	return b.buildWrapper(buildfn)
}

func (b *componentBuilderBase) buildVolume(_ int32) componentBuilder {
	return b.buildWrapper(nil)
}

func (b *componentBuilderBase) buildVolumeMount(idx int32) componentBuilder {
	buildfn := func() ([]client.Object, error) {
		if b.concreteBuilder.mutableWorkload(idx) == nil {
			return nil, fmt.Errorf("build volume mount but workload is nil, cluster: %s, component: %s",
				b.Comp.GetClusterName(), b.Comp.GetName())
		}

		podSpec := b.concreteBuilder.mutablePodSpec(idx)
		for _, cc := range []*[]corev1.Container{&podSpec.Containers, &podSpec.InitContainers} {
			volumes := podSpec.Volumes
			for _, c := range *cc {
				for _, v := range c.VolumeMounts {
					// if persistence is not found, add emptyDir pod.spec.volumes[]
					createfn := func(_ string) corev1.Volume {
						return corev1.Volume{
							Name: v.Name,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						}
					}
					volumes, _ = intctrlutil.CreateOrUpdateVolume(volumes, v.Name, createfn, nil)
				}
			}
			podSpec.Volumes = volumes
		}
		return nil, nil
	}
	return b.buildWrapper(buildfn)
}

func (b *componentBuilderBase) buildTLSCert() componentBuilder {
	buildfn := func() ([]client.Object, error) {
		cluster := b.Comp.GetCluster()
		component := b.Comp.GetSynthesizedComponent()
		if !component.TLS {
			return nil, nil
		}
		if component.Issuer == nil {
			return nil, fmt.Errorf("issuer shouldn't be nil when tls enabled")
		}

		objs := make([]client.Object, 0)
		switch component.Issuer.Name {
		case appsv1alpha1.IssuerUserProvided:
			if err := plan.CheckTLSSecretRef(b.ReqCtx, b.Client, cluster.Namespace, component.Issuer.SecretRef); err != nil {
				return nil, err
			}
		case appsv1alpha1.IssuerKubeBlocks:
			secret, err := plan.ComposeTLSSecret(cluster.Namespace, cluster.Name, component.Name)
			if err != nil {
				return nil, err
			}
			objs = append(objs, secret)
		}
		return objs, nil
	}
	return b.buildWrapper(buildfn)
}

func (b *componentBuilderBase) buildTLSVolume(idx int32) componentBuilder {
	buildfn := func() ([]client.Object, error) {
		if b.concreteBuilder.mutableWorkload(idx) == nil {
			return nil, fmt.Errorf("build TLS volumes but workload is nil, cluster: %s, component: %s",
				b.Comp.GetClusterName(), b.Comp.GetName())
		}
		// build secret volume and volume mount
		podSpec := b.concreteBuilder.mutablePodSpec(idx)
		return nil, updateTLSVolumeAndVolumeMount(podSpec, b.Comp.GetClusterName(), *b.Comp.GetSynthesizedComponent())
	}
	return b.buildWrapper(buildfn)
}

func (b *componentBuilderBase) complete() error {
	if b.Error != nil {
		return b.Error
	}
	workload := b.concreteBuilder.mutableWorkload(0)
	if workload == nil {
		return fmt.Errorf("fail to create compoennt workloads, cluster: %s, component: %s",
			b.Comp.GetClusterName(), b.Comp.GetName())
	}
	b.Comp.addWorkload(workload, b.defaultAction, nil)
	return nil
}

func (b *componentBuilderBase) buildWrapper(buildfn func() ([]client.Object, error)) componentBuilder {
	if b.Error != nil || buildfn == nil {
		return b.concreteBuilder
	}
	objs, err := buildfn()
	if err != nil {
		b.Error = err
	} else {
		for _, obj := range objs {
			b.Comp.addResource(obj, b.defaultAction, nil)
		}
	}
	return b.concreteBuilder
}
