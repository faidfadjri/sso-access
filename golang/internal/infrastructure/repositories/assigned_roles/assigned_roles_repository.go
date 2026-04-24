package assigned_roles

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
	"log"

	"gorm.io/gorm"
)

type assignedRolesRepository struct {
	db *gorm.DB
}

func NewAssignedRolesRepository(db *gorm.DB) AssignedRolesRepository {
	return &assignedRolesRepository{
		db: db,
	}
}

func (r *assignedRolesRepository) Create(usr *domain.UserServiceRole) error {
	var existing domain.UserServiceRole
	err := r.db.Unscoped().Where("user_id = ? AND role_id = ? AND service_id = ?", usr.UserId, usr.RoleId, usr.ServiceId).First(&existing).Error

	if err == nil {
		// If it exists (even soft-deleted), restore it by setting deleted_at to null
		return r.db.Unscoped().Model(&existing).Update("deleted_at", nil).Error
	} else if err == gorm.ErrRecordNotFound {
		// If it does not exist, create it anew
		return r.db.Create(usr).Error
	}

	return err
}

func (r *assignedRolesRepository) GetList(p *pagination.Pagination) (*pagination.Pagination, error) {
	var list []entities.AssignedRoles

	// Use domain.UserServiceRole{} for the model so GORM knows the correct table
	db := r.db.Model(&domain.UserServiceRole{}).
		Select("users.*, sr.service_role_id, sr.role_name, s.service_id, s.service_name").
		Joins("LEFT JOIN users ON users.user_id = user_service_role.user_id").
		Joins("LEFT JOIN service_roles sr ON sr.service_role_id = user_service_role.role_id").
		Joins("LEFT JOIN services s ON s.service_id = sr.service_id")

	
	if p.Search != "" {
		db = db.Where("users.full_name LIKE ? OR users.email LIKE ? OR users.username LIKE ? OR users.phone LIKE ? OR s.service_name LIKE ?", "%"+p.Search+"%", "%"+p.Search+"%", "%"+p.Search+"%", "%"+p.Search+"%", "%"+p.Search+"%")
	}

	if p.ServiceId != nil {
		db = db.Where("s.service_id = ?", *p.ServiceId)
	}

	if p.RoleId != nil {
		db = db.Where("sr.service_role_id = ?", *p.RoleId)
	}

	// Pass the model definition to Paginate so it counts from the correct table
	err := db.Scopes(p.Paginate(&domain.UserServiceRole{}, db)).
		Find(&list).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	p.Rows = list
	return p, nil
}

func (r *assignedRolesRepository) GetByKeys(userID, roleID uint64) (*entities.AssignedRoles, error) {
	var usr entities.AssignedRoles
	
	err := r.db.Model(&domain.UserServiceRole{}).
		Select("users.*, sr.service_role_id, sr.role_name, s.service_id, s.service_name").
		Joins("LEFT JOIN users ON users.user_id = user_service_role.user_id").
		Joins("LEFT JOIN service_roles sr ON sr.service_role_id = user_service_role.role_id").
		Joins("LEFT JOIN services s ON s.service_id = sr.service_id").
		Where("user_service_role.user_id = ? AND user_service_role.role_id = ?", userID, roleID).
		First(&usr).Error

	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func (r *assignedRolesRepository) GetServiceIdByRoleId(roleID uint64) (uint64, error) {
	var serviceRoles domain.ServiceRoles
	
	err := r.db.Model(&domain.ServiceRoles{}).
		Select("service_id").
		Where("service_role_id = ?", roleID).
		First(&serviceRoles).Error

	if err != nil {
		return 0, err
	}
	return serviceRoles.ServiceId, nil
}

func (r *assignedRolesRepository) Delete(userID, roleID uint64) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&domain.UserServiceRole{}).Error
}

func (r *assignedRolesRepository) DeleteByUserAndService(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.UserServiceRole{}).Error
}
