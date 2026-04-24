package role

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) Create(role *domain.ServiceRoles) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetRoles(p *pagination.Pagination) (*pagination.Pagination, error) {
	var roles []entities.Roles

	db := r.db.Model(&domain.ServiceRoles{}).
		Select("service_roles.*, s.service_name").
		Joins("LEFT JOIN services s ON s.service_id = service_roles.service_id")

	if p.Search != "" {
		db = db.Where("service_roles.role_name LIKE ?", "%"+p.Search+"%")
	}

	err := db.Scopes(p.Paginate(&domain.ServiceRoles{}, db)).Scan(&roles).Error
	if err != nil {
		return nil, err
	}

	p.Rows = roles
	return p, nil
}

func (r *roleRepository) GetRoleByID(id uint64) (*domain.ServiceRoles, error) {
	var role domain.ServiceRoles
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(role *domain.ServiceRoles) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.ServiceRoles{}, id).Error
}
