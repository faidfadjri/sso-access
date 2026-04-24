package assigned_roles

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
)

type AssignedRolesRepository interface {
	Create(usr *domain.UserServiceRole) error
	GetList(p *pagination.Pagination) (*pagination.Pagination, error)
	GetByKeys(userID, roleID uint64) (*entities.AssignedRoles, error)
	GetServiceIdByRoleId(roleID uint64) (uint64, error)
	Delete(userID, roleID uint64) error
	DeleteByUserAndService(userID uint64) error
}