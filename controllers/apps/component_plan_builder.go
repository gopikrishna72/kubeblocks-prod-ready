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
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/constant"
	roclient "github.com/apecloud/kubeblocks/pkg/controller/client"
	"github.com/apecloud/kubeblocks/pkg/controller/component"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
	intctrlutil "github.com/apecloud/kubeblocks/pkg/controllerutil"
)

// componentTransformContext a graph.TransformContext implementation for Cluster reconciliation
type componentTransformContext struct {
	context.Context
	Client roclient.ReadonlyClient
	record.EventRecorder
	logr.Logger
	Cluster             *appsv1alpha1.Cluster
	CompDef             *appsv1alpha1.ComponentDefinition
	Component           *appsv1alpha1.Component
	ComponentOrig       *appsv1alpha1.Component
	SynthesizeComponent *component.SynthesizedComponent
}

func (c *componentTransformContext) GetContext() context.Context {
	return c.Context
}

func (c *componentTransformContext) GetClient() roclient.ReadonlyClient {
	return c.Client
}

func (c *componentTransformContext) GetRecorder() record.EventRecorder {
	return c.EventRecorder
}

func (c *componentTransformContext) GetLogger() logr.Logger {
	return c.Logger
}

// componentPlanBuilder a graph.PlanBuilder implementation for Component reconciliation
type componentPlanBuilder struct {
	req          ctrl.Request
	cli          client.Client
	transCtx     *componentTransformContext
	transformers graph.TransformerChain
}

// clusterPlan a graph.Plan implementation for Cluster reconciliation
type componentPlan struct {
	dag      *graph.DAG
	walkFunc graph.WalkFunc
	cli      client.Client
	transCtx *componentTransformContext
}

var _ graph.TransformContext = &componentTransformContext{}
var _ graph.PlanBuilder = &componentPlanBuilder{}
var _ graph.Plan = &componentPlan{}

// PlanBuilder implementation

func (c *componentPlanBuilder) Init() error {
	comp := &appsv1alpha1.Component{}
	if err := c.cli.Get(c.transCtx.Context, c.req.NamespacedName, comp); err != nil {
		return err
	}

	c.transCtx.Component = comp
	c.transCtx.ComponentOrig = comp.DeepCopy()
	c.transformers = append(c.transformers, &componentInitTransformer{
		Component:     c.transCtx.Component,
		ComponentOrig: c.transCtx.ComponentOrig,
	})
	return nil
}

func (c *componentPlanBuilder) AddTransformer(transformer ...graph.Transformer) graph.PlanBuilder {
	c.transformers = append(c.transformers, transformer...)
	return c
}

func (c *componentPlanBuilder) AddParallelTransformer(transformer ...graph.Transformer) graph.PlanBuilder {
	c.transformers = append(c.transformers, &ParallelTransformers{transformers: transformer})
	return c
}

// Build runs all transformers to generate a plan
func (c *componentPlanBuilder) Build() (graph.Plan, error) {
	var err error
	// new a DAG and apply chain on it
	dag := graph.NewDAG()
	err = c.transformers.ApplyTo(c.transCtx, dag)
	c.transCtx.Logger.V(1).Info(fmt.Sprintf("DAG: %s", dag))

	// construct execution plan
	plan := &componentPlan{
		dag:      dag,
		walkFunc: c.componentWalkFunc,
		cli:      c.cli,
		transCtx: c.transCtx,
	}
	return plan, err
}

// Plan implementation

func (p *componentPlan) Execute() error {
	return p.dag.WalkReverseTopoOrder(p.walkFunc, nil)
}

// Do the real works

// NewComponentPlanBuilder returns a componentPlanBuilder powered PlanBuilder
func NewComponentPlanBuilder(ctx intctrlutil.RequestCtx, cli client.Client, req ctrl.Request) graph.PlanBuilder {
	return &componentPlanBuilder{
		req: req,
		cli: cli,
		transCtx: &componentTransformContext{
			Context:       ctx.Ctx,
			Client:        model.NewGraphClient(cli),
			EventRecorder: ctx.Recorder,
			Logger:        ctx.Log,
		},
	}
}

func (c *componentPlanBuilder) componentWalkFunc(v graph.Vertex) error {
	vertex, ok := v.(*model.ObjectVertex)
	if !ok {
		return fmt.Errorf("wrong vertex type %v", v)
	}
	if vertex.Action == nil {
		return errors.New("vertex action can't be nil")
	}
	switch *vertex.Action {
	case model.CREATE:
		err := c.cli.Create(c.transCtx.Context, vertex.Obj)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	case model.UPDATE:
		err := c.cli.Update(c.transCtx.Context, vertex.Obj)
		if err != nil && !apierrors.IsNotFound(err) {
			c.transCtx.Logger.Error(err, fmt.Sprintf("update %T error: %s", vertex.Obj, vertex.OriObj.GetName()))
			return err
		}
	case model.DELETE:
		if controllerutil.RemoveFinalizer(vertex.Obj, constant.DBComponentFinalizerName) {
			err := c.cli.Update(c.transCtx.Context, vertex.Obj)
			if err != nil && !apierrors.IsNotFound(err) {
				c.transCtx.Logger.Error(err, fmt.Sprintf("delete %T error: %s", vertex.Obj, vertex.Obj.GetName()))
				return err
			}
		}
		if !model.IsObjectDeleting(vertex.Obj) {
			err := c.cli.Delete(c.transCtx.Context, vertex.Obj)
			if err != nil && !apierrors.IsNotFound(err) {
				return err
			}
		}
	case model.STATUS:
		if err := c.cli.Status().Update(c.transCtx.Context, vertex.Obj); err != nil {
			return err
		}
	}
	return nil
}
