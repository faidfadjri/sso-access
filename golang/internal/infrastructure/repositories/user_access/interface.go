package user_access

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
)

type UserAccessRepository interface {
	Create(access *domain.UserAccess) error
	GetUserAccesses(p *pagination.Pagination) (*pagination.Pagination, error)
	GetUserAccessByID(id uint64) (*domain.UserAccess, error)
	GetUserAccessByUserId(userId uint64) ([]domain.UserAccess, error)
	IsUserAccessExist(userId uint64, serviceId uint64) (bool, error)
	Update(access *domain.UserAccess) error
	Delete(id uint64) error	
}
