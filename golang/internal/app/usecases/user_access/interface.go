package user_access

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
)

type UserAccessUsecase interface {
	CreateUserAccess(userID uint64, serviceIDs []uint64, status string) ([]*domain.UserAccess, error)
	GetUserAccesses(p *pagination.Pagination) (*pagination.Pagination, error)
	GetUserAccessByID(id uint64) (*domain.UserAccess, error)
	UpdateUserAccess(userID uint64, serviceIDs []uint64, status string) error
	DeleteUserAccess(id uint64) error
}