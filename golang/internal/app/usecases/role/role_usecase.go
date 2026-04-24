package role

import (
	"akastra-access/internal/infrastructure/databases/domain"
	repoRole "akastra-access/internal/infrastructure/repositories/role"
	"akastra-access/internal/infrastructure/utils/pagination"
	pkgErrors "akastra-access/internal/pkg/errors"
	"errors"

	"gorm.io/gorm"
)

type roleUsecase struct {
	roleRepo repoRole.RoleRepository
}

func NewRoleUsecase(repo repoRole.RoleRepository) RoleUsecase {
	return &roleUsecase{
		roleRepo: repo,
	}
}

func (u *roleUsecase) CreateRole(roleName string, serviceID uint64) (*domain.ServiceRoles, error) {
	role := &domain.ServiceRoles{
		RoleName:  roleName,
		ServiceId: serviceID,
	}

	if err := u.roleRepo.Create(role); err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}

	return role, nil
}

func (u *roleUsecase) GetRoles(p *pagination.Pagination) (*pagination.Pagination, error) {
	roles, err := u.roleRepo.GetRoles(p)
	if err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}
	return roles, nil
}

func (u *roleUsecase) GetRoleByID(id uint64) (*domain.ServiceRoles, error) {
	role, err := u.roleRepo.GetRoleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, pkgErrors.ErrFailedQuery
	}
	return role, nil
}

func (u *roleUsecase) UpdateRole(id uint64, roleName string, serviceID uint64) error {
	role, err := u.roleRepo.GetRoleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	role.RoleName = roleName
	role.ServiceId = serviceID

	if err := u.roleRepo.Update(role); err != nil {
		return pkgErrors.ErrFailedQuery
	}

	return nil
}

func (u *roleUsecase) DeleteRole(id uint64) error {
	if _, err := u.roleRepo.GetRoleByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	if err := u.roleRepo.Delete(id); err != nil {
		return pkgErrors.ErrFailedDelete
	}

	return nil
}
