package mocks

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"github.com/stretchr/testify/mock"
)

type RolePermissionRepositoryMock struct {
	mock.Mock
}

func (m *RolePermissionRepositoryMock) Create(rp *domain.RolePermissions) error {
	args := m.Called(rp)
	return args.Error(0)
}

func (m *RolePermissionRepositoryMock) GetList(p *pagination.Pagination) (*pagination.Pagination, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Pagination), args.Error(1)
}

func (m *RolePermissionRepositoryMock) GetByKeys(roleID, permissionID uint64) (*domain.RolePermissions, error) {
	args := m.Called(roleID, permissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RolePermissions), args.Error(1)
}

func (m *RolePermissionRepositoryMock) Delete(roleID, permissionID uint64) error {
	args := m.Called(roleID, permissionID)
	return args.Error(0)
}

func (m *RolePermissionRepositoryMock) DeleteByRole(roleID uint64) error {
	args := m.Called(roleID)
	return args.Error(0)
}
