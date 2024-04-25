// Code generated by MockGen. DO NOT EDIT.
// Source: ./bash.go

// Package mock_util is a generated GoMock package.
package mock_util

import (
	bytes "bytes"
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIBashUtil is a mock of IBashUtil interface.
type MockIBashUtil struct {
	ctrl     *gomock.Controller
	recorder *MockIBashUtilMockRecorder
}

// MockIBashUtilMockRecorder is the mock recorder for MockIBashUtil.
type MockIBashUtilMockRecorder struct {
	mock *MockIBashUtil
}

// NewMockIBashUtil creates a new mock instance.
func NewMockIBashUtil(ctrl *gomock.Controller) *MockIBashUtil {
	mock := &MockIBashUtil{ctrl: ctrl}
	mock.recorder = &MockIBashUtilMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBashUtil) EXPECT() *MockIBashUtilMockRecorder {
	return m.recorder
}

// GetBashFileBody mocks base method.
func (m *MockIBashUtil) GetBashFileBody(arg0 *multipart.FileHeader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBashFileBody", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBashFileBody indicates an expected call of GetBashFileBody.
func (mr *MockIBashUtilMockRecorder) GetBashFileBody(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBashFileBody", reflect.TypeOf((*MockIBashUtil)(nil).GetBashFileBody), arg0)
}

// GetBashFileBuffer mocks base method.
func (m *MockIBashUtil) GetBashFileBuffer(arg0 string) *bytes.Buffer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBashFileBuffer", arg0)
	ret0, _ := ret[0].(*bytes.Buffer)
	return ret0
}

// GetBashFileBuffer indicates an expected call of GetBashFileBuffer.
func (mr *MockIBashUtilMockRecorder) GetBashFileBuffer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBashFileBuffer", reflect.TypeOf((*MockIBashUtil)(nil).GetBashFileBuffer), arg0)
}

// GetBashFileExtension mocks base method.
func (m *MockIBashUtil) GetBashFileExtension(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBashFileExtension", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetBashFileExtension indicates an expected call of GetBashFileExtension.
func (mr *MockIBashUtilMockRecorder) GetBashFileExtension(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBashFileExtension", reflect.TypeOf((*MockIBashUtil)(nil).GetBashFileExtension), arg0)
}

// GetBashFileTitle mocks base method.
func (m *MockIBashUtil) GetBashFileTitle(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBashFileTitle", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetBashFileTitle indicates an expected call of GetBashFileTitle.
func (mr *MockIBashUtilMockRecorder) GetBashFileTitle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBashFileTitle", reflect.TypeOf((*MockIBashUtil)(nil).GetBashFileTitle), arg0)
}

// ValidateBashFileExtension mocks base method.
func (m *MockIBashUtil) ValidateBashFileExtension(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateBashFileExtension", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ValidateBashFileExtension indicates an expected call of ValidateBashFileExtension.
func (mr *MockIBashUtilMockRecorder) ValidateBashFileExtension(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateBashFileExtension", reflect.TypeOf((*MockIBashUtil)(nil).ValidateBashFileExtension), arg0)
}
