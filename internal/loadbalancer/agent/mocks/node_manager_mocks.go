/*
Copyright 2022 The KubeBlocks Authors

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

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/apecloud/kubeblocks/internal/loadbalancer/agent (interfaces: NodeManager)

// Package mock_agent is a generated GoMock package.
package mock_agent

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	agent "github.com/apecloud/kubeblocks/internal/loadbalancer/agent"
)

// MockNodeManager is a mock of NodeManager interface.
type MockNodeManager struct {
	ctrl     *gomock.Controller
	recorder *MockNodeManagerMockRecorder
}

// MockNodeManagerMockRecorder is the mock recorder for MockNodeManager.
type MockNodeManagerMockRecorder struct {
	mock *MockNodeManager
}

// NewMockNodeManager creates a new mock instance.
func NewMockNodeManager(ctrl *gomock.Controller) *MockNodeManager {
	mock := &MockNodeManager{ctrl: ctrl}
	mock.recorder = &MockNodeManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeManager) EXPECT() *MockNodeManagerMockRecorder {
	return m.recorder
}

// ChooseSpareNode mocks base method.
func (m *MockNodeManager) ChooseSpareNode(arg0 string) (agent.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChooseSpareNode", arg0)
	ret0, _ := ret[0].(agent.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChooseSpareNode indicates an expected call of ChooseSpareNode.
func (mr *MockNodeManagerMockRecorder) ChooseSpareNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChooseSpareNode", reflect.TypeOf((*MockNodeManager)(nil).ChooseSpareNode), arg0)
}

// GetNode mocks base method.
func (m *MockNodeManager) GetNode(arg0 string) (agent.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNode", arg0)
	ret0, _ := ret[0].(agent.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNode indicates an expected call of GetNode.
func (mr *MockNodeManagerMockRecorder) GetNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockNodeManager)(nil).GetNode), arg0)
}

// GetNodes mocks base method.
func (m *MockNodeManager) GetNodes() ([]agent.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodes")
	ret0, _ := ret[0].([]agent.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes.
func (mr *MockNodeManagerMockRecorder) GetNodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockNodeManager)(nil).GetNodes))
}
