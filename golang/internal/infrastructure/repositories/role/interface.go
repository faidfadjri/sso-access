package role

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
)

type RoleRepository interface {
	Create(role *domain.ServiceRoles) error
	GetRoles(p *pagination.Pagination) (*pagination.Pagination, error)
	GetRoleByID(id uint64) (*domain.ServiceRoles, error)
	Update(role *domain.ServiceRoles) error
	Delete(id uint64) error
}	