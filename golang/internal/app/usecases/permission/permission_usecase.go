package permission

import (
	"akastra-access/internal/infrastructure/databases/domain"
	repoPermission "akastra-access/internal/infrastructure/repositories/permission"
	"akastra-access/internal/infrastructure/utils/pagination"
	pkgErrors "akastra-access/internal/pkg/errors"
	"errors"

	"gorm.io/gorm"
)

type PermissionUsecase interface {
	CreatePermission(key string, description string) (*domain.Permissions, error)
	GetPermissions(p *pagination.Pagination) (*pagination.Pagination, error)
	GetPermissionByID(id uint64) (*domain.Permissions, error)
	UpdatePermission(id uint64, key string, description string) error
	DeletePermission(id uint64) error
}

type permissionUsecase struct {
	permRepo repoPermission.PermissionRepository
}

func NewPermissionUsecase(repo repoPermission.PermissionRepository) PermissionUsecase {
	return &permissionUsecase{
		permRepo: repo,
	}
}

func (u *permissionUsecase) CreatePermission(key string, description string) (*domain.Permissions, error) {
	permission := &domain.Permissions{
		PermissionKey: key,
		Description:   description,
	}

	if err := u.permRepo.Create(permission); err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}

	return permission, nil
}

func (u *permissionUsecase) GetPermissions(p *pagination.Pagination) (*pagination.Pagination, error) {
	permissions, err := u.permRepo.GetPermissions(p)
	if err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}
	return permissions, nil
}

func (u *permissionUsecase) GetPermissionByID(id uint64) (*domain.Permissions, error) {
	permission, err := u.permRepo.GetPermissionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, pkgErrors.ErrFailedQuery
	}
	return permission, nil
}

func (u *permissionUsecase) UpdatePermission(id uint64, key string, description string) error {
	permission, err := u.permRepo.GetPermissionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	permission.PermissionKey = key
	permission.Description = description

	if err := u.permRepo.Update(permission); err != nil {
		return pkgErrors.ErrFailedQuery
	}

	return nil
}

func (u *permissionUsecase) DeletePermission(id uint64) error {
	if _, err := u.permRepo.GetPermissionByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	if err := u.permRepo.Delete(id); err != nil {
		return pkgErrors.ErrFailedDelete
	}

	return nil
}
