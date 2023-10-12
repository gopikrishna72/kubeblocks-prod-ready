/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/apecloud/kubeblocks/apis/dataprotection/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeActionSets implements ActionSetInterface
type FakeActionSets struct {
	Fake *FakeDataprotectionV1alpha1
}

var actionsetsResource = v1alpha1.SchemeGroupVersion.WithResource("actionsets")

var actionsetsKind = v1alpha1.SchemeGroupVersion.WithKind("ActionSet")

// Get takes name of the actionSet, and returns the corresponding actionSet object, and an error if there is any.
func (c *FakeActionSets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ActionSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(actionsetsResource, name), &v1alpha1.ActionSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActionSet), err
}

// List takes label and field selectors, and returns the list of ActionSets that match those selectors.
func (c *FakeActionSets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ActionSetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(actionsetsResource, actionsetsKind, opts), &v1alpha1.ActionSetList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ActionSetList{ListMeta: obj.(*v1alpha1.ActionSetList).ListMeta}
	for _, item := range obj.(*v1alpha1.ActionSetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested actionSets.
func (c *FakeActionSets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(actionsetsResource, opts))
}

// Create takes the representation of a actionSet and creates it.  Returns the server's representation of the actionSet, and an error, if there is any.
func (c *FakeActionSets) Create(ctx context.Context, actionSet *v1alpha1.ActionSet, opts v1.CreateOptions) (result *v1alpha1.ActionSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(actionsetsResource, actionSet), &v1alpha1.ActionSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActionSet), err
}

// Update takes the representation of a actionSet and updates it. Returns the server's representation of the actionSet, and an error, if there is any.
func (c *FakeActionSets) Update(ctx context.Context, actionSet *v1alpha1.ActionSet, opts v1.UpdateOptions) (result *v1alpha1.ActionSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(actionsetsResource, actionSet), &v1alpha1.ActionSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActionSet), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeActionSets) UpdateStatus(ctx context.Context, actionSet *v1alpha1.ActionSet, opts v1.UpdateOptions) (*v1alpha1.ActionSet, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(actionsetsResource, "status", actionSet), &v1alpha1.ActionSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActionSet), err
}

// Delete takes name of the actionSet and deletes it. Returns an error if one occurs.
func (c *FakeActionSets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(actionsetsResource, name, opts), &v1alpha1.ActionSet{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeActionSets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(actionsetsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ActionSetList{})
	return err
}

// Patch applies the patch and returns the patched actionSet.
func (c *FakeActionSets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ActionSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(actionsetsResource, name, pt, data, subresources...), &v1alpha1.ActionSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActionSet), err
}
