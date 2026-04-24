package permission

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(permission *domain.Permissions) error
	GetPermissions(p *pagination.Pagination) (*pagination.Pagination, error)
	GetPermissionByID(id uint64) (*domain.Permissions, error)
	Update(permission *domain.Permissions) error
	Delete(id uint64) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) Create(permission *domain.Permissions) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) GetPermissions(p *pagination.Pagination) (*pagination.Pagination, error) {
	var permissions []domain.Permissions

	db := r.db.Model(&domain.Permissions{})

	if p.Search != "" {
		db = db.Where("permission_key LIKE ? OR description LIKE ?", "%"+p.Search+"%", "%"+p.Search+"%")
	}

	err := db.Scopes(p.Paginate(&permissions, db)).Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	p.Rows = permissions
	return p, nil
}

func (r *permissionRepository) GetPermissionByID(id uint64) (*domain.Permissions, error) {
	var permission domain.Permissions
	if err := r.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Update(permission *domain.Permissions) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Permissions{}, id).Error
}
