// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/controller-runtime/pkg/manager (interfaces: Manager)
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=controllers/mocks/manager_mock.go sigs.k8s.io/controller-runtime/pkg/manager Manager
//

// Package mocks is a generated GoMock package.
package mocks

/*
Copyright 2021-2025 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	context "context"
	http "net/http"
	reflect "reflect"

	logr "github.com/go-logr/logr"
	gomock "go.uber.org/mock/gomock"
	meta "k8s.io/apimachinery/pkg/api/meta"
	runtime "k8s.io/apimachinery/pkg/runtime"
	rest "k8s.io/client-go/rest"
	record "k8s.io/client-go/tools/record"
	cache "sigs.k8s.io/controller-runtime/pkg/cache"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	config "sigs.k8s.io/controller-runtime/pkg/config"
	healthz "sigs.k8s.io/controller-runtime/pkg/healthz"
	manager "sigs.k8s.io/controller-runtime/pkg/manager"
	webhook "sigs.k8s.io/controller-runtime/pkg/webhook"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockManager) Add(arg0 manager.Runnable) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockManagerMockRecorder) Add(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockManager)(nil).Add), arg0)
}

// AddHealthzCheck mocks base method.
func (m *MockManager) AddHealthzCheck(arg0 string, arg1 healthz.Checker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddHealthzCheck", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddHealthzCheck indicates an expected call of AddHealthzCheck.
func (mr *MockManagerMockRecorder) AddHealthzCheck(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHealthzCheck", reflect.TypeOf((*MockManager)(nil).AddHealthzCheck), arg0, arg1)
}

// AddMetricsServerExtraHandler mocks base method.
func (m *MockManager) AddMetricsServerExtraHandler(arg0 string, arg1 http.Handler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetricsServerExtraHandler", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMetricsServerExtraHandler indicates an expected call of AddMetricsServerExtraHandler.
func (mr *MockManagerMockRecorder) AddMetricsServerExtraHandler(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetricsServerExtraHandler", reflect.TypeOf((*MockManager)(nil).AddMetricsServerExtraHandler), arg0, arg1)
}

// AddReadyzCheck mocks base method.
func (m *MockManager) AddReadyzCheck(arg0 string, arg1 healthz.Checker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReadyzCheck", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddReadyzCheck indicates an expected call of AddReadyzCheck.
func (mr *MockManagerMockRecorder) AddReadyzCheck(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReadyzCheck", reflect.TypeOf((*MockManager)(nil).AddReadyzCheck), arg0, arg1)
}

// Elected mocks base method.
func (m *MockManager) Elected() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Elected")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Elected indicates an expected call of Elected.
func (mr *MockManagerMockRecorder) Elected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Elected", reflect.TypeOf((*MockManager)(nil).Elected))
}

// GetAPIReader mocks base method.
func (m *MockManager) GetAPIReader() client.Reader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAPIReader")
	ret0, _ := ret[0].(client.Reader)
	return ret0
}

// GetAPIReader indicates an expected call of GetAPIReader.
func (mr *MockManagerMockRecorder) GetAPIReader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIReader", reflect.TypeOf((*MockManager)(nil).GetAPIReader))
}

// GetCache mocks base method.
func (m *MockManager) GetCache() cache.Cache {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCache")
	ret0, _ := ret[0].(cache.Cache)
	return ret0
}

// GetCache indicates an expected call of GetCache.
func (mr *MockManagerMockRecorder) GetCache() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCache", reflect.TypeOf((*MockManager)(nil).GetCache))
}

// GetClient mocks base method.
func (m *MockManager) GetClient() client.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClient")
	ret0, _ := ret[0].(client.Client)
	return ret0
}

// GetClient indicates an expected call of GetClient.
func (mr *MockManagerMockRecorder) GetClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClient", reflect.TypeOf((*MockManager)(nil).GetClient))
}

// GetConfig mocks base method.
func (m *MockManager) GetConfig() *rest.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig")
	ret0, _ := ret[0].(*rest.Config)
	return ret0
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockManagerMockRecorder) GetConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockManager)(nil).GetConfig))
}

// GetControllerOptions mocks base method.
func (m *MockManager) GetControllerOptions() config.Controller {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetControllerOptions")
	ret0, _ := ret[0].(config.Controller)
	return ret0
}

// GetControllerOptions indicates an expected call of GetControllerOptions.
func (mr *MockManagerMockRecorder) GetControllerOptions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetControllerOptions", reflect.TypeOf((*MockManager)(nil).GetControllerOptions))
}

// GetEventRecorderFor mocks base method.
func (m *MockManager) GetEventRecorderFor(arg0 string) record.EventRecorder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventRecorderFor", arg0)
	ret0, _ := ret[0].(record.EventRecorder)
	return ret0
}

// GetEventRecorderFor indicates an expected call of GetEventRecorderFor.
func (mr *MockManagerMockRecorder) GetEventRecorderFor(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventRecorderFor", reflect.TypeOf((*MockManager)(nil).GetEventRecorderFor), arg0)
}

// GetFieldIndexer mocks base method.
func (m *MockManager) GetFieldIndexer() client.FieldIndexer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFieldIndexer")
	ret0, _ := ret[0].(client.FieldIndexer)
	return ret0
}

// GetFieldIndexer indicates an expected call of GetFieldIndexer.
func (mr *MockManagerMockRecorder) GetFieldIndexer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFieldIndexer", reflect.TypeOf((*MockManager)(nil).GetFieldIndexer))
}

// GetHTTPClient mocks base method.
func (m *MockManager) GetHTTPClient() *http.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHTTPClient")
	ret0, _ := ret[0].(*http.Client)
	return ret0
}

// GetHTTPClient indicates an expected call of GetHTTPClient.
func (mr *MockManagerMockRecorder) GetHTTPClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHTTPClient", reflect.TypeOf((*MockManager)(nil).GetHTTPClient))
}

// GetLogger mocks base method.
func (m *MockManager) GetLogger() logr.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogger")
	ret0, _ := ret[0].(logr.Logger)
	return ret0
}

// GetLogger indicates an expected call of GetLogger.
func (mr *MockManagerMockRecorder) GetLogger() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogger", reflect.TypeOf((*MockManager)(nil).GetLogger))
}

// GetRESTMapper mocks base method.
func (m *MockManager) GetRESTMapper() meta.RESTMapper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRESTMapper")
	ret0, _ := ret[0].(meta.RESTMapper)
	return ret0
}

// GetRESTMapper indicates an expected call of GetRESTMapper.
func (mr *MockManagerMockRecorder) GetRESTMapper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRESTMapper", reflect.TypeOf((*MockManager)(nil).GetRESTMapper))
}

// GetScheme mocks base method.
func (m *MockManager) GetScheme() *runtime.Scheme {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScheme")
	ret0, _ := ret[0].(*runtime.Scheme)
	return ret0
}

// GetScheme indicates an expected call of GetScheme.
func (mr *MockManagerMockRecorder) GetScheme() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScheme", reflect.TypeOf((*MockManager)(nil).GetScheme))
}

// GetWebhookServer mocks base method.
func (m *MockManager) GetWebhookServer() webhook.Server {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWebhookServer")
	ret0, _ := ret[0].(webhook.Server)
	return ret0
}

// GetWebhookServer indicates an expected call of GetWebhookServer.
func (mr *MockManagerMockRecorder) GetWebhookServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWebhookServer", reflect.TypeOf((*MockManager)(nil).GetWebhookServer))
}

// Start mocks base method.
func (m *MockManager) Start(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockManagerMockRecorder) Start(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockManager)(nil).Start), arg0)
}
