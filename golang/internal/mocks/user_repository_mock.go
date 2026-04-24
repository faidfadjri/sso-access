package mocks

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user *domain.Users) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetUsers(p *pagination.Pagination) (*pagination.Pagination, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Pagination), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByID(id uint64) (*domain.Users, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Users), args.Error(1)
}

func (m *UserRepositoryMock) Update(user *domain.Users) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UserRepositoryMock) BatchDelete(ids []uint64) error {
	args := m.Called(ids)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetUserByEmail(email string) (*domain.Users, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Users), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByUsername(username string) (*domain.Users, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Users), args.Error(1)
}
