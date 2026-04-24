package mocks

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"github.com/stretchr/testify/mock"
)

type RoleRepositoryMock struct {
	mock.Mock
}

func (m *RoleRepositoryMock) Create(role *domain.ServiceRoles) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *RoleRepositoryMock) GetRoles(p *pagination.Pagination) (*pagination.Pagination, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Pagination), args.Error(1)
}

func (m *RoleRepositoryMock) GetRoleByID(id uint64) (*domain.ServiceRoles, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ServiceRoles), args.Error(1)
}

func (m *RoleRepositoryMock) Update(role *domain.ServiceRoles) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *RoleRepositoryMock) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}
