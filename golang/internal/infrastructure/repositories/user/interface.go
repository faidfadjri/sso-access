package user

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
)

type UserRepository interface {
	Create(user *domain.Users) error
	GetUsers(p *pagination.Pagination) (*pagination.Pagination, error)
	GetUserByID(id uint64) (*domain.Users, error)
	Update(user *domain.Users) error
	Delete(id uint64) error
	BatchDelete(ids []uint64) error
	
	// Helper to check existing
	GetUserByEmail(email string) (*domain.Users, error)
	GetUserByUsername(username string) (*domain.Users, error)
}
