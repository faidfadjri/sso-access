package user_access

import (
	"akastra-access/internal/app/config"
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
	"log"

	"gorm.io/gorm"
)

type userAccessRepository struct {
	db *gorm.DB
}

func NewUserAccessRepository(db *gorm.DB) UserAccessRepository {
	return &userAccessRepository{
		db: db,
	}
}

func (r *userAccessRepository) Create(access *domain.UserAccess) error {
	return r.db.Create(access).Error
}

func (r *userAccessRepository) GetUserAccesses(p *pagination.Pagination) (*pagination.Pagination, error) {
	var accesses []entities.UserAccess
	frontendURL := config.Load().FrontendURL

	db := r.db.Model(&domain.UserAccess{}).
		Select(`
			user_access.*, 
			users.full_name, 
			users.email, 
			users.username, 
			users.phone, 
			services.redirect_url, 
			services.service_name,
			services.client_id,
			DATE_FORMAT(user_access.created_at, '%Y-%m-%d %H:%i') as created_at
		`).
		Joins("LEFT JOIN users ON users.user_id = user_access.user_id").
		Joins("LEFT JOIN services ON services.service_id = user_access.service_id").
		Where("services.is_active != 0").Where("services.redirect_url NOT LIKE ?", "%"+ frontendURL +"%")

	if p.Search != "" {
		db = db.Where("users.full_name LIKE ? OR users.email LIKE ? OR users.username LIKE ? OR users.phone LIKE ? OR services.service_name LIKE ?", "%"+p.Search+"%", "%"+p.Search+"%", "%"+p.Search+"%", "%"+p.Search+"%", "%"+p.Search+"%")
	}

	if p.UserId != nil {
		db = db.Where("user_access.user_id = ?", *p.UserId)
	}

	if p.ServiceId != nil {
		db = db.Where("user_access.service_id = ?", *p.ServiceId)
	}

	if p.RoleId != nil {
		db = db.Where("user_access.role_id = ?", *p.RoleId)
	}

	err := db.Scopes(p.Paginate(&domain.UserAccess{}, db)).Scan(&accesses).Error
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	p.Rows = accesses
	return p, nil
}

func (r *userAccessRepository) GetUserAccessByID(id uint64) (*domain.UserAccess, error) {
	var access domain.UserAccess
	if err := r.db.First(&access, id).Error; err != nil {
		return nil, err
	}
	return &access, nil
}

func (r *userAccessRepository) Update(access *domain.UserAccess) error {
	return r.db.Save(access).Error
}

func (r *userAccessRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.UserAccess{}, id).Error
}

func (r *userAccessRepository) GetUserAccessByUserId(userId uint64) ([]domain.UserAccess, error) {
	var accesses []domain.UserAccess
	if err := r.db.Where("user_id = ?", userId).Find(&accesses).Error; err != nil {
		return nil, err
	}
	return accesses, nil
}

func (r *userAccessRepository) IsUserAccessExist(userId uint64, serviceId uint64) (bool, error) {
	var access domain.UserAccess
	if err := r.db.Where("user_id = ? AND service_id = ?", userId, serviceId).First(&access).Error; err != nil {
		return false, err
	}
	return true, nil
}
