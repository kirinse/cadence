// The MIT License (MIT)

// Copyright (c) 2017-2020 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Code generated by MockGen. DO NOT EDIT.
// Source: resolver.go

// Package membership is a generated GoMock package.
package membership

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockResolver is a mock of Resolver interface.
type MockResolver struct {
	ctrl     *gomock.Controller
	recorder *MockResolverMockRecorder
}

// MockResolverMockRecorder is the mock recorder for MockResolver.
type MockResolverMockRecorder struct {
	mock *MockResolver
}

// NewMockResolver creates a new mock instance.
func NewMockResolver(ctrl *gomock.Controller) *MockResolver {
	mock := &MockResolver{ctrl: ctrl}
	mock.recorder = &MockResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResolver) EXPECT() *MockResolverMockRecorder {
	return m.recorder
}

// EvictSelf mocks base method.
func (m *MockResolver) EvictSelf() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EvictSelf")
	ret0, _ := ret[0].(error)
	return ret0
}

// EvictSelf indicates an expected call of EvictSelf.
func (mr *MockResolverMockRecorder) EvictSelf() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EvictSelf", reflect.TypeOf((*MockResolver)(nil).EvictSelf))
}

// Lookup mocks base method.
func (m *MockResolver) Lookup(service, key string) (*HostInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lookup", service, key)
	ret0, _ := ret[0].(*HostInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lookup indicates an expected call of Lookup.
func (mr *MockResolverMockRecorder) Lookup(service, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lookup", reflect.TypeOf((*MockResolver)(nil).Lookup), service, key)
}

// MemberCount mocks base method.
func (m *MockResolver) MemberCount(service string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MemberCount", service)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MemberCount indicates an expected call of MemberCount.
func (mr *MockResolverMockRecorder) MemberCount(service interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MemberCount", reflect.TypeOf((*MockResolver)(nil).MemberCount), service)
}

// Members mocks base method.
func (m *MockResolver) Members(service string) ([]*HostInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Members", service)
	ret0, _ := ret[0].([]*HostInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Members indicates an expected call of Members.
func (mr *MockResolverMockRecorder) Members(service interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Members", reflect.TypeOf((*MockResolver)(nil).Members), service)
}

// Start mocks base method.
func (m *MockResolver) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockResolverMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockResolver)(nil).Start))
}

// Stop mocks base method.
func (m *MockResolver) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockResolverMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockResolver)(nil).Stop))
}

// Subscribe mocks base method.
func (m *MockResolver) Subscribe(service, name string, notifyChannel chan<- *ChangedEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", service, name, notifyChannel)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockResolverMockRecorder) Subscribe(service, name, notifyChannel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockResolver)(nil).Subscribe), service, name, notifyChannel)
}

// Unsubscribe mocks base method.
func (m *MockResolver) Unsubscribe(service, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", service, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockResolverMockRecorder) Unsubscribe(service, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockResolver)(nil).Unsubscribe), service, name)
}

// WhoAmI mocks base method.
func (m *MockResolver) WhoAmI() (*HostInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WhoAmI")
	ret0, _ := ret[0].(*HostInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WhoAmI indicates an expected call of WhoAmI.
func (mr *MockResolverMockRecorder) WhoAmI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WhoAmI", reflect.TypeOf((*MockResolver)(nil).WhoAmI))
}
