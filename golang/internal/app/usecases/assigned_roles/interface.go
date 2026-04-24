package assigned_roles

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
)

type AssignedRolesUsecase interface {
	AssignRole(userIDs []uint64, roleID uint64) ([]*domain.UserServiceRole, error)
	GetList(p *pagination.Pagination) (*pagination.Pagination, error)
	GetByKeys(userID, roleID uint64) (*entities.AssignedRoles, error)
	RemoveRole(userID, roleID uint64) error
}
