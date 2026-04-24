package user

import (
	"akastra-access/internal/app/config"
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *domain.Users) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUsers(p *pagination.Pagination) (*pagination.Pagination, error) {
	var users []entities.UserWithService

	serviceID := config.Load().IdentityProviderID

	db := r.db.Model(&domain.Users{}).
		Select("users.*, sr.role_name, s.service_name").
		Joins("LEFT JOIN user_service_role usr ON usr.user_id = users.user_id AND usr.service_id = ?", serviceID).
		Joins("LEFT JOIN services s ON s.service_id = usr.service_id").
		Joins("LEFT JOIN service_roles sr ON sr.service_role_id = usr.role_id");

	if p.Search != "" {
		db = db.Where("full_name LIKE ? OR email LIKE ? OR username LIKE ?", "%"+p.Search+"%", "%"+p.Search+"%", "%"+p.Search+"%")
	}
		

	err := db.Scopes(p.Paginate(&domain.Users{}, db)).Scan(&users).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	p.Rows = users
	return p, nil
}

func (r *userRepository) GetUserByID(id uint64) (*domain.Users, error) {
	var user domain.Users
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.Users) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Users{}, id).Error
}

func (r *userRepository) BatchDelete(ids []uint64) error {
	return r.db.Delete(&domain.Users{}, ids).Error
}

func (r *userRepository) GetUserByEmail(email string) (*domain.Users, error) {
	var user domain.Users
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*domain.Users, error) {
	var user domain.Users
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
