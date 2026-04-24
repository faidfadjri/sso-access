package mocks

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"

	"github.com/stretchr/testify/mock"
)

type OAuthRepositoryMock struct {
	mock.Mock
}

func (m *OAuthRepositoryMock) GetUserByEmailorUsername(emailOrUsername string) (*domain.Users, error) {
	args := m.Called(emailOrUsername)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Users), args.Error(1)
}

func (m *OAuthRepositoryMock) CreateUser(req *domain.Users) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *OAuthRepositoryMock) GetUserByID(userID uint64) (*domain.Users, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Users), args.Error(1)
}

func (m *OAuthRepositoryMock) GetUserByIDWithService(userID uint64, serviceID uint64) (*entities.UserWithService, error) {
	args := m.Called(userID, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserWithService), args.Error(1)
}

func (m *OAuthRepositoryMock) UpdateUser(req *domain.Users) error {
	args := m.Called(req)
	return args.Error(0)
}
