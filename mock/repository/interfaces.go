// Code generated by MockGen. DO NOT EDIT.
// Source: repository/interfaces.go

// Package repository_mock is a generated GoMock package.
package repository_mock

import (
	context "context"
	reflect "reflect"
	repository "swpr/repository"

	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetTestById mocks base method.
func (m *MockRepositoryInterface) GetTestById(ctx context.Context, input repository.GetTestByIdInput) (repository.GetTestByIdOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTestById", ctx, input)
	ret0, _ := ret[0].(repository.GetTestByIdOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTestById indicates an expected call of GetTestById.
func (mr *MockRepositoryInterfaceMockRecorder) GetTestById(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTestById", reflect.TypeOf((*MockRepositoryInterface)(nil).GetTestById), ctx, input)
}

// LoginAttemptCreate mocks base method.
func (m *MockRepositoryInterface) LoginAttemptCreate(ctx context.Context, input repository.LoginAttemptCreate) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginAttemptCreate", ctx, input)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginAttemptCreate indicates an expected call of LoginAttemptCreate.
func (mr *MockRepositoryInterfaceMockRecorder) LoginAttemptCreate(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginAttemptCreate", reflect.TypeOf((*MockRepositoryInterface)(nil).LoginAttemptCreate), ctx, input)
}

// UserCreate mocks base method.
func (m *MockRepositoryInterface) UserCreate(ctx context.Context, input repository.UserCreate) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserCreate", ctx, input)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserCreate indicates an expected call of UserCreate.
func (mr *MockRepositoryInterfaceMockRecorder) UserCreate(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserCreate", reflect.TypeOf((*MockRepositoryInterface)(nil).UserCreate), ctx, input)
}

// UserGetById mocks base method.
func (m *MockRepositoryInterface) UserGetById(ctx context.Context, id int64) (*repository.UserGet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserGetById", ctx, id)
	ret0, _ := ret[0].(*repository.UserGet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserGetById indicates an expected call of UserGetById.
func (mr *MockRepositoryInterfaceMockRecorder) UserGetById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserGetById", reflect.TypeOf((*MockRepositoryInterface)(nil).UserGetById), ctx, id)
}

// UserGetByPhone mocks base method.
func (m *MockRepositoryInterface) UserGetByPhone(ctx context.Context, phone string) (*repository.UserGet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserGetByPhone", ctx, phone)
	ret0, _ := ret[0].(*repository.UserGet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserGetByPhone indicates an expected call of UserGetByPhone.
func (mr *MockRepositoryInterfaceMockRecorder) UserGetByPhone(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserGetByPhone", reflect.TypeOf((*MockRepositoryInterface)(nil).UserGetByPhone), ctx, phone)
}

// UserUpdate mocks base method.
func (m *MockRepositoryInterface) UserUpdate(ctx context.Context, input repository.UserUpdate) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserUpdate", ctx, input)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserUpdate indicates an expected call of UserUpdate.
func (mr *MockRepositoryInterfaceMockRecorder) UserUpdate(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserUpdate", reflect.TypeOf((*MockRepositoryInterface)(nil).UserUpdate), ctx, input)
}
