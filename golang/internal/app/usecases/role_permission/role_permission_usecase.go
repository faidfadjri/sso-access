package role_permission

import (
	"akastra-access/internal/infrastructure/databases/domain"
	repoRP "akastra-access/internal/infrastructure/repositories/role_permission"
	"akastra-access/internal/infrastructure/utils/pagination"
	pkgErrors "akastra-access/internal/pkg/errors"
	"errors"

	"gorm.io/gorm"
)

type RolePermissionUsecase interface {
	AssignPermission(roleID, permissionID uint64) (*domain.RolePermissions, error)
	GetList(p *pagination.Pagination) (*pagination.Pagination, error)
	GetByKeys(roleID, permissionID uint64) (*domain.RolePermissions, error)
	RemovePermission(roleID, permissionID uint64) error
}

type rolePermissionUsecase struct {
	repo repoRP.RolePermissionRepository
}

func NewRolePermissionUsecase(repo repoRP.RolePermissionRepository) RolePermissionUsecase {
	return &rolePermissionUsecase{
		repo: repo,
	}
}

func (u *rolePermissionUsecase) AssignPermission(roleID, permissionID uint64) (*domain.RolePermissions, error) {
	rp := &domain.RolePermissions{
		RoleId:       roleID,
		PermissionId: permissionID,
	}

	if err := u.repo.Create(rp); err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}

	return rp, nil
}

func (u *rolePermissionUsecase) GetList(p *pagination.Pagination) (*pagination.Pagination, error) {
	list, err := u.repo.GetList(p)
	if err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}
	return list, nil
}

func (u *rolePermissionUsecase) GetByKeys(roleID, permissionID uint64) (*domain.RolePermissions, error) {
	rp, err := u.repo.GetByKeys(roleID, permissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, pkgErrors.ErrFailedQuery
	}
	return rp, nil
}

func (u *rolePermissionUsecase) RemovePermission(roleID, permissionID uint64) error {
	if _, err := u.repo.GetByKeys(roleID, permissionID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	if err := u.repo.Delete(roleID, permissionID); err != nil {
		return pkgErrors.ErrFailedDelete
	}

	return nil
}
