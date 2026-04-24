package mocks

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"github.com/stretchr/testify/mock"
)

type ServiceRepositoryMock struct {
	mock.Mock
}

func (m *ServiceRepositoryMock) CreateClient(req *domain.Services) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *ServiceRepositoryMock) GetClientByIds(serviceIds []uint64) ([]*domain.Services, error) {
	args := m.Called(serviceIds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Services), args.Error(1)
}

func (m *ServiceRepositoryMock) GetClientById(serviceId uint64) (*domain.Services, error) {
	args := m.Called(serviceId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Services), args.Error(1)
}

func (m *ServiceRepositoryMock) GetClientByClientId(clientId string) (*domain.Services, error) {
	args := m.Called(clientId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Services), args.Error(1)
}

func (m *ServiceRepositoryMock) DeleteClientById(serviceId uint64) error {
	args := m.Called(serviceId)
	return args.Error(0)
}

func (m *ServiceRepositoryMock) GetClients(p *pagination.Pagination) (*pagination.Pagination, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Pagination), args.Error(1)
}

func (m *ServiceRepositoryMock) UpdateClient(req *domain.Services) error {
	args := m.Called(req)
	return args.Error(0)
}
