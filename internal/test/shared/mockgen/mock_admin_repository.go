// Code generated by MockGen. DO NOT EDIT.
// Source: ./admin_repository.go

// Package mockgen is a generated GoMock package.
package mockgen

import (
	reflect "reflect"

	domain "github.com/CamiloLeonP/parking-radar/internal/app/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockIAdminRepository is a mock of IAdminRepository interface.
type MockIAdminRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAdminRepositoryMockRecorder
}

// MockIAdminRepositoryMockRecorder is the mock recorder for MockIAdminRepository.
type MockIAdminRepositoryMockRecorder struct {
	mock *MockIAdminRepository
}

// NewMockIAdminRepository creates a new mock instance.
func NewMockIAdminRepository(ctrl *gomock.Controller) *MockIAdminRepository {
	mock := &MockIAdminRepository{ctrl: ctrl}
	mock.recorder = &MockIAdminRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAdminRepository) EXPECT() *MockIAdminRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIAdminRepository) Create(admin *domain.Admin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", admin)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIAdminRepositoryMockRecorder) Create(admin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIAdminRepository)(nil).Create), admin)
}

// ExistsByAuth0UUID mocks base method.
func (m *MockIAdminRepository) ExistsByAuth0UUID(auth0UUID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByAuth0UUID", auth0UUID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExistsByAuth0UUID indicates an expected call of ExistsByAuth0UUID.
func (mr *MockIAdminRepositoryMockRecorder) ExistsByAuth0UUID(auth0UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByAuth0UUID", reflect.TypeOf((*MockIAdminRepository)(nil).ExistsByAuth0UUID), auth0UUID)
}

// FindByAuth0UUID mocks base method.
func (m *MockIAdminRepository) FindByAuth0UUID(auth0UUID string) (*domain.Admin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByAuth0UUID", auth0UUID)
	ret0, _ := ret[0].(*domain.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByAuth0UUID indicates an expected call of FindByAuth0UUID.
func (mr *MockIAdminRepositoryMockRecorder) FindByAuth0UUID(auth0UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByAuth0UUID", reflect.TypeOf((*MockIAdminRepository)(nil).FindByAuth0UUID), auth0UUID)
}

// Update mocks base method.
func (m *MockIAdminRepository) Update(admin *domain.Admin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", admin)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIAdminRepositoryMockRecorder) Update(admin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIAdminRepository)(nil).Update), admin)
}
