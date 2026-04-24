package service

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
)

// Service defines the interface for Service use case
type ServiceRepository interface {
	CreateClient(req *domain.Services) error
	GetClientById(serviceId uint64) (*domain.Services, error)
	GetClientByIds(serviceIds []uint64) ([]*domain.Services, error)

	GetClientByClientId(clientId string) (*domain.Services, error)
	DeleteClientById(serviceId uint64) error
	GetClients(pagination *pagination.Pagination) (*pagination.Pagination, error)
	UpdateClient(service *domain.Services) error
}
