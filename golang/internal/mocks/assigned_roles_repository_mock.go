package mocks

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"github.com/stretchr/testify/mock"
)

type AssignedRolesRepositoryMock struct {
	mock.Mock
}

func (m *AssignedRolesRepositoryMock) Create(usr *domain.UserServiceRole) error {
	args := m.Called(usr)
	return args.Error(0)
}

func (m *AssignedRolesRepositoryMock) GetList(p *pagination.Pagination) (*pagination.Pagination, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Pagination), args.Error(1)
}

func (m *AssignedRolesRepositoryMock) GetByKeys(userID, roleID uint64) (*entities.AssignedRoles, error) {
	args := m.Called(userID, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AssignedRoles), args.Error(1)
}

func (m *AssignedRolesRepositoryMock) Delete(userID, roleID uint64) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *AssignedRolesRepositoryMock) DeleteByUserAndService(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *AssignedRolesRepositoryMock) GetServiceIdByRoleId(roleID uint64) (uint64, error) {
	args := m.Called(roleID)
	return args.Get(0).(uint64), args.Error(1)
}
