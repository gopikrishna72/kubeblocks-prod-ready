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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/apecloud/kubeblocks/internal/constant"
	"github.com/apecloud/kubeblocks/internal/controller/graph"
)

type ComponentAssureMetaTransformer struct{}

var _ graph.Transformer = &ComponentAssureMetaTransformer{}

func (t *ComponentAssureMetaTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*ComponentTransformContext)
	component := transCtx.Component

	// The object is not being deleted, so if it does not have our finalizer,
	// then lets add the finalizer and update the object. This is equivalent
	// registering our finalizer.
	if !controllerutil.ContainsFinalizer(component, constant.DBComponentFinalizerName) {
		controllerutil.AddFinalizer(component, constant.DBComponentFinalizerName)
	}

	// patch the label to prevent the label from being modified by the user.
	labels := component.Labels
	if labels == nil {
		labels = map[string]string{}
	}
	labelName := labels[constant.ComponentDefinitionLabelKey]
	if labelName != component.Spec.CompDef {
		labels[constant.ComponentDefinitionLabelKey] = component.Spec.CompDef
		component.Labels = labels
	}
	return nil
}