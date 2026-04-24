package mocks

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"github.com/stretchr/testify/mock"
)

type UserServiceRoleRepositoryMock struct {
	mock.Mock
}

func (m *UserServiceRoleRepositoryMock) Create(usr *domain.UserServiceRole) error {
	args := m.Called(usr)
	return args.Error(0)
}

func (m *UserServiceRoleRepositoryMock) GetList(p *pagination.Pagination) (*pagination.Pagination, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Pagination), args.Error(1)
}

func (m *UserServiceRoleRepositoryMock) GetByKeys(userID, roleID uint64) (*domain.UserServiceRole, error) {
	args := m.Called(userID, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserServiceRole), args.Error(1)
}

func (m *UserServiceRoleRepositoryMock) Delete(userID, roleID uint64) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *UserServiceRoleRepositoryMock) DeleteByUserAndService(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}
