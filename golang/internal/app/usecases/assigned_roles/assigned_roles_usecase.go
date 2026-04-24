package assigned_roles

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
	repoUSR "akastra-access/internal/infrastructure/repositories/assigned_roles"
	"akastra-access/internal/infrastructure/utils/pagination"
	pkgErrors "akastra-access/internal/pkg/errors"
	"errors"
	"log"
	"strings"

	"gorm.io/gorm"
)


type assignedRolesUsecase struct {
	repo repoUSR.AssignedRolesRepository
}

func NewAssignedRolesUsecase(repo repoUSR.AssignedRolesRepository) AssignedRolesUsecase {
	return &assignedRolesUsecase{
		repo: repo,
	}
}

func (u *assignedRolesUsecase) AssignRole(userIDs []uint64, roleID uint64) ([]*domain.UserServiceRole, error) {
	var createdRoles []*domain.UserServiceRole
	serviceID, err := u.repo.GetServiceIdByRoleId(roleID)
	if err != nil {
		return nil, pkgErrors.ErrInvalidRoleId
	}

	for _, userID := range userIDs {
		usr := &domain.UserServiceRole{
			UserId: userID,
			RoleId: roleID,
			ServiceId: serviceID,
		}

		if err := u.repo.Create(usr); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				continue
			}
			log.Println(err)
			return nil, pkgErrors.ErrFailedQuery
		}
		createdRoles = append(createdRoles, usr)
	}

	return createdRoles, nil
}

func (u *assignedRolesUsecase) GetList(p *pagination.Pagination) (*pagination.Pagination, error) {
	list, err := u.repo.GetList(p)
	if err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}
	return list, nil
}

func (u *assignedRolesUsecase) GetByKeys(userID, roleID uint64) (*entities.AssignedRoles, error) {
	usr, err := u.repo.GetByKeys(userID, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, pkgErrors.ErrFailedQuery
	}
	return usr, nil
}

func (u *assignedRolesUsecase) RemoveRole(userID, roleID uint64) error {
	if _, err := u.repo.GetByKeys(userID, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	if err := u.repo.Delete(userID, roleID); err != nil {
		return pkgErrors.ErrFailedDelete
	}

	return nil
}
