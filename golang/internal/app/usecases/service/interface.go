package service

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"mime/multipart"
)

// Service defines the interface for Service use case
type ServiceUsecase interface {
	CreateClient(req request.CreateClientReq, file multipart.File, header *multipart.FileHeader) (*domain.Services, error)
	DeleteClientById(serviceId uint64) error
	GetClients(pagination *pagination.Pagination) (*pagination.Pagination, error)
	UpdateClient(req request.UpdateClientReq, file multipart.File, header *multipart.FileHeader) error
}
