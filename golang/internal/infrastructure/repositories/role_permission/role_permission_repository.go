package role_permission

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"gorm.io/gorm"
)

type RolePermissionRepository interface {
	Create(rp *domain.RolePermissions) error
	GetList(p *pagination.Pagination) (*pagination.Pagination, error)
	GetByKeys(roleID, permissionID uint64) (*domain.RolePermissions, error)
	Delete(roleID, permissionID uint64) error
	DeleteByRole(roleID uint64) error
}

type rolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) RolePermissionRepository {
	return &rolePermissionRepository{
		db: db,
	}
}

func (r *rolePermissionRepository) Create(rp *domain.RolePermissions) error {
	return r.db.Create(rp).Error
}

func (r *rolePermissionRepository) GetList(p *pagination.Pagination) (*pagination.Pagination, error) {
	var list []domain.RolePermissions

	db := r.db.Model(&domain.RolePermissions{})

	err := db.Scopes(p.Paginate(&list, db)).Find(&list).Error
	if err != nil {
		return nil, err
	}

	p.Rows = list
	return p, nil
}

func (r *rolePermissionRepository) GetByKeys(roleID, permissionID uint64) (*domain.RolePermissions, error) {
	var rp domain.RolePermissions
	err := r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&rp).Error
	if err != nil {
		return nil, err
	}
	return &rp, nil
}

func (r *rolePermissionRepository) Delete(roleID, permissionID uint64) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&domain.RolePermissions{}).Error
}

func (r *rolePermissionRepository) DeleteByRole(roleID uint64) error {
	return r.db.Where("role_id = ?", roleID).Delete(&domain.RolePermissions{}).Error
}
