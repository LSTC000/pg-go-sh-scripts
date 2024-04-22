// Code generated by MockGen. DO NOT EDIT.
// Source: ./bash.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	dto "pg-sh-scripts/internal/dto"
	model "pg-sh-scripts/internal/model"
	alias "pg-sh-scripts/internal/type/alias"
	pagination "pg-sh-scripts/pkg/sql/pagination"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
)

// MockIBashService is a mock of IBashService interface.
type MockIBashService struct {
	ctrl     *gomock.Controller
	recorder *MockIBashServiceMockRecorder
}

// MockIBashServiceMockRecorder is the mock recorder for MockIBashService.
type MockIBashServiceMockRecorder struct {
	mock *MockIBashService
}

// NewMockIBashService creates a new mock instance.
func NewMockIBashService(ctrl *gomock.Controller) *MockIBashService {
	mock := &MockIBashService{ctrl: ctrl}
	mock.recorder = &MockIBashServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBashService) EXPECT() *MockIBashServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIBashService) Create(ctx context.Context, dto dto.CreateBash) (*model.Bash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, dto)
	ret0, _ := ret[0].(*model.Bash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIBashServiceMockRecorder) Create(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIBashService)(nil).Create), ctx, dto)
}

// GetOneById mocks base method.
func (m *MockIBashService) GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneById", ctx, id)
	ret0, _ := ret[0].(*model.Bash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOneById indicates an expected call of GetOneById.
func (mr *MockIBashServiceMockRecorder) GetOneById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneById", reflect.TypeOf((*MockIBashService)(nil).GetOneById), ctx, id)
}

// GetPaginationPage mocks base method.
func (m *MockIBashService) GetPaginationPage(ctx context.Context, paginationParams pagination.LimitOffsetParams) (alias.BashLimitOffsetPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaginationPage", ctx, paginationParams)
	ret0, _ := ret[0].(alias.BashLimitOffsetPage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPaginationPage indicates an expected call of GetPaginationPage.
func (mr *MockIBashServiceMockRecorder) GetPaginationPage(ctx, paginationParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaginationPage", reflect.TypeOf((*MockIBashService)(nil).GetPaginationPage), ctx, paginationParams)
}

// RemoveById mocks base method.
func (m *MockIBashService) RemoveById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveById", ctx, id)
	ret0, _ := ret[0].(*model.Bash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveById indicates an expected call of RemoveById.
func (mr *MockIBashServiceMockRecorder) RemoveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveById", reflect.TypeOf((*MockIBashService)(nil).RemoveById), ctx, id)
}
