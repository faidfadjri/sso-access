package oauth

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"

	"gorm.io/gorm"
)

type oauthRepository struct {
	db *gorm.DB
}

func NewOAuthRepository(db *gorm.DB) *oauthRepository {
	return &oauthRepository{db: db}
}

func (r *oauthRepository) GetUserByEmailorUsername(emailOrUsername string) (*domain.Users, error) {
	var user domain.Users
	if emailOrUsername != "" {
		if err := r.db.Where("email = ?", emailOrUsername).Or("username = ?", emailOrUsername).First(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

func (r *oauthRepository) CreateUser(req *domain.Users) error {
	return r.db.Create(req).Error
}

func (r *oauthRepository) GetUserByID(userID uint64) (*domain.Users, error) {
	var user domain.Users
	if err := r.db.Where("user_id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *oauthRepository) GetUserByIDWithService(userID uint64, serviceID uint64) (*entities.UserWithService, error) {
	var user entities.UserWithService

	err := r.db.
		Model(&domain.Users{}).
		Select(`
			users.user_id,
			users.full_name,
			users.email,
			users.username,
			users.photo,
			users.phone,
			s.service_name,
			sr.role_name
		`).
		Joins(`
			LEFT JOIN user_service_role usr 
				ON usr.user_id = users.user_id 
				AND usr.service_id = ?
		`, serviceID).
		Joins(`
			LEFT JOIN services s 
				ON s.service_id = usr.service_id
		`).
		Joins(`
			LEFT JOIN service_roles sr 
				ON sr.service_role_id = usr.role_id
		`).
		Where("users.user_id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}


func (r *oauthRepository) UpdateUser(req *domain.Users) error {
	return r.db.Save(req).Error
}
