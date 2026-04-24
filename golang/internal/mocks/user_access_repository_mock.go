package mocks

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"github.com/stretchr/testify/mock"
)

type UserAccessRepositoryMock struct {
	mock.Mock
}

func (m *UserAccessRepositoryMock) Create(access *domain.UserAccess) error {
	args := m.Called(access)
	return args.Error(0)
}

func (m *UserAccessRepositoryMock) GetUserAccesses(p *pagination.Pagination) (*pagination.Pagination, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Pagination), args.Error(1)
}

func (m *UserAccessRepositoryMock) GetUserAccessByID(id uint64) (*domain.UserAccess, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserAccess), args.Error(1)
}

func (m *UserAccessRepositoryMock) GetUserAccessByUserId(userId uint64) ([]domain.UserAccess, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserAccess), args.Error(1)
}

func (m *UserAccessRepositoryMock) IsUserAccessExist(userId uint64, serviceId uint64) (bool, error) {
	args := m.Called(userId, serviceId)
	return args.Bool(0), args.Error(1)
}

func (m *UserAccessRepositoryMock) Update(access *domain.UserAccess) error {
	args := m.Called(access)
	return args.Error(0)
}

func (m *UserAccessRepositoryMock) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}
