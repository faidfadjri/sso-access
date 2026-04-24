package role

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
)

type RoleUsecase interface {
	CreateRole(roleName string, serviceID uint64) (*domain.ServiceRoles, error)
	GetRoles(p *pagination.Pagination) (*pagination.Pagination, error)
	GetRoleByID(id uint64) (*domain.ServiceRoles, error)
	UpdateRole(id uint64, roleName string, serviceID uint64) error
	DeleteRole(id uint64) error
}
