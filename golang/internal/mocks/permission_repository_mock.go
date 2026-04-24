package mocks

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"github.com/stretchr/testify/mock"
)

type PermissionRepositoryMock struct {
	mock.Mock
}

func (m *PermissionRepositoryMock) Create(permission *domain.Permissions) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *PermissionRepositoryMock) GetPermissions(p *pagination.Pagination) (*pagination.Pagination, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Pagination), args.Error(1)
}

func (m *PermissionRepositoryMock) GetPermissionByID(id uint64) (*domain.Permissions, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permissions), args.Error(1)
}

func (m *PermissionRepositoryMock) Update(permission *domain.Permissions) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *PermissionRepositoryMock) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}
