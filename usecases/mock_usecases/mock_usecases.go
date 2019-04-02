// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/int128/goxzst/usecases/interfaces (interfaces: Make,CrossBuild,CreateZip,CreateSHA,RenderTemplate)

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	gomock "github.com/golang/mock/gomock"
	interfaces "github.com/int128/goxzst/usecases/interfaces"
	reflect "reflect"
)

// MockMake is a mock of Make interface
type MockMake struct {
	ctrl     *gomock.Controller
	recorder *MockMakeMockRecorder
}

// MockMakeMockRecorder is the mock recorder for MockMake
type MockMakeMockRecorder struct {
	mock *MockMake
}

// NewMockMake creates a new mock instance
func NewMockMake(ctrl *gomock.Controller) *MockMake {
	mock := &MockMake{ctrl: ctrl}
	mock.recorder = &MockMakeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMake) EXPECT() *MockMakeMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockMake) Do(arg0 interfaces.MakeIn) error {
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do
func (mr *MockMakeMockRecorder) Do(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockMake)(nil).Do), arg0)
}

// MockCrossBuild is a mock of CrossBuild interface
type MockCrossBuild struct {
	ctrl     *gomock.Controller
	recorder *MockCrossBuildMockRecorder
}

// MockCrossBuildMockRecorder is the mock recorder for MockCrossBuild
type MockCrossBuildMockRecorder struct {
	mock *MockCrossBuild
}

// NewMockCrossBuild creates a new mock instance
func NewMockCrossBuild(ctrl *gomock.Controller) *MockCrossBuild {
	mock := &MockCrossBuild{ctrl: ctrl}
	mock.recorder = &MockCrossBuildMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCrossBuild) EXPECT() *MockCrossBuildMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockCrossBuild) Do(arg0 interfaces.CrossBuildIn) error {
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do
func (mr *MockCrossBuildMockRecorder) Do(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockCrossBuild)(nil).Do), arg0)
}

// MockCreateZip is a mock of CreateZip interface
type MockCreateZip struct {
	ctrl     *gomock.Controller
	recorder *MockCreateZipMockRecorder
}

// MockCreateZipMockRecorder is the mock recorder for MockCreateZip
type MockCreateZipMockRecorder struct {
	mock *MockCreateZip
}

// NewMockCreateZip creates a new mock instance
func NewMockCreateZip(ctrl *gomock.Controller) *MockCreateZip {
	mock := &MockCreateZip{ctrl: ctrl}
	mock.recorder = &MockCreateZipMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCreateZip) EXPECT() *MockCreateZipMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockCreateZip) Do(arg0 interfaces.CreateZipIn) error {
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do
func (mr *MockCreateZipMockRecorder) Do(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockCreateZip)(nil).Do), arg0)
}

// MockCreateSHA is a mock of CreateSHA interface
type MockCreateSHA struct {
	ctrl     *gomock.Controller
	recorder *MockCreateSHAMockRecorder
}

// MockCreateSHAMockRecorder is the mock recorder for MockCreateSHA
type MockCreateSHAMockRecorder struct {
	mock *MockCreateSHA
}

// NewMockCreateSHA creates a new mock instance
func NewMockCreateSHA(ctrl *gomock.Controller) *MockCreateSHA {
	mock := &MockCreateSHA{ctrl: ctrl}
	mock.recorder = &MockCreateSHAMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCreateSHA) EXPECT() *MockCreateSHAMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockCreateSHA) Do(arg0 interfaces.CreateSHAIn) (*interfaces.CreateSHAOut, error) {
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*interfaces.CreateSHAOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockCreateSHAMockRecorder) Do(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockCreateSHA)(nil).Do), arg0)
}

// MockRenderTemplate is a mock of RenderTemplate interface
type MockRenderTemplate struct {
	ctrl     *gomock.Controller
	recorder *MockRenderTemplateMockRecorder
}

// MockRenderTemplateMockRecorder is the mock recorder for MockRenderTemplate
type MockRenderTemplateMockRecorder struct {
	mock *MockRenderTemplate
}

// NewMockRenderTemplate creates a new mock instance
func NewMockRenderTemplate(ctrl *gomock.Controller) *MockRenderTemplate {
	mock := &MockRenderTemplate{ctrl: ctrl}
	mock.recorder = &MockRenderTemplateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRenderTemplate) EXPECT() *MockRenderTemplateMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockRenderTemplate) Do(arg0 interfaces.RenderTemplateIn) error {
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do
func (mr *MockRenderTemplateMockRecorder) Do(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockRenderTemplate)(nil).Do), arg0)
}
