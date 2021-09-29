/*
Copyright 2021 The k8gb Contributors.

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
// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/infobloxopen/infoblox-go-client (interfaces: IBConnector)

// Package dns is a generated GoMock package.
package dns

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

// MockIBConnector is a mock of IBConnector interface.
type MockIBConnector struct {
	ctrl     *gomock.Controller
	recorder *MockIBConnectorMockRecorder
}

// MockIBConnectorMockRecorder is the mock recorder for MockIBConnector.
type MockIBConnectorMockRecorder struct {
	mock *MockIBConnector
}

// NewMockIBConnector creates a new mock instance.
func NewMockIBConnector(ctrl *gomock.Controller) *MockIBConnector {
	mock := &MockIBConnector{ctrl: ctrl}
	mock.recorder = &MockIBConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBConnector) EXPECT() *MockIBConnectorMockRecorder {
	return m.recorder
}

// CreateObject mocks base method.
func (m *MockIBConnector) CreateObject(arg0 ibclient.IBObject) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateObject", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateObject indicates an expected call of CreateObject.
func (mr *MockIBConnectorMockRecorder) CreateObject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateObject", reflect.TypeOf((*MockIBConnector)(nil).CreateObject), arg0)
}

// DeleteObject mocks base method.
func (m *MockIBConnector) DeleteObject(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObject", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteObject indicates an expected call of DeleteObject.
func (mr *MockIBConnectorMockRecorder) DeleteObject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObject", reflect.TypeOf((*MockIBConnector)(nil).DeleteObject), arg0)
}

// GetObject mocks base method.
func (m *MockIBConnector) GetObject(arg0 ibclient.IBObject, arg1 string, arg2 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObject", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetObject indicates an expected call of GetObject.
func (mr *MockIBConnectorMockRecorder) GetObject(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObject", reflect.TypeOf((*MockIBConnector)(nil).GetObject), arg0, arg1, arg2)
}

// UpdateObject mocks base method.
func (m *MockIBConnector) UpdateObject(arg0 ibclient.IBObject, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateObject", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateObject indicates an expected call of UpdateObject.
func (mr *MockIBConnectorMockRecorder) UpdateObject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateObject", reflect.TypeOf((*MockIBConnector)(nil).UpdateObject), arg0, arg1)
}
