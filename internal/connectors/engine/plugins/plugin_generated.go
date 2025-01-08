// Code generated by MockGen. DO NOT EDIT.
// Source: plugin.go
//
// Generated by this command:
//
//	mockgen -source plugin.go -destination plugin_generated.go -package plugins . Plugins
//

// Package plugins is a generated GoMock package.
package plugins

import (
	json "encoding/json"
	reflect "reflect"

	models "github.com/formancehq/payments/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockPlugins is a mock of Plugins interface.
type MockPlugins struct {
	ctrl     *gomock.Controller
	recorder *MockPluginsMockRecorder
}

// MockPluginsMockRecorder is the mock recorder for MockPlugins.
type MockPluginsMockRecorder struct {
	mock *MockPlugins
}

// NewMockPlugins creates a new mock instance.
func NewMockPlugins(ctrl *gomock.Controller) *MockPlugins {
	mock := &MockPlugins{ctrl: ctrl}
	mock.recorder = &MockPluginsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlugins) EXPECT() *MockPluginsMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPlugins) Get(arg0 models.ConnectorID) (models.Plugin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(models.Plugin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPluginsMockRecorder) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPlugins)(nil).Get), arg0)
}

// GetConfig mocks base method.
func (m *MockPlugins) GetConfig(arg0 models.ConnectorID) (models.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig", arg0)
	ret0, _ := ret[0].(models.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockPluginsMockRecorder) GetConfig(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockPlugins)(nil).GetConfig), arg0)
}

// RegisterPlugin mocks base method.
func (m *MockPlugins) RegisterPlugin(arg0 models.ConnectorID, arg1 string, arg2 models.Config, arg3 json.RawMessage, arg4 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterPlugin", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterPlugin indicates an expected call of RegisterPlugin.
func (mr *MockPluginsMockRecorder) RegisterPlugin(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPlugin", reflect.TypeOf((*MockPlugins)(nil).RegisterPlugin), arg0, arg1, arg2, arg3, arg4)
}

// UnregisterPlugin mocks base method.
func (m *MockPlugins) UnregisterPlugin(arg0 models.ConnectorID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterPlugin", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterPlugin indicates an expected call of UnregisterPlugin.
func (mr *MockPluginsMockRecorder) UnregisterPlugin(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterPlugin", reflect.TypeOf((*MockPlugins)(nil).UnregisterPlugin), arg0)
}
