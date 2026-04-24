package oauth

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/infrastructure/databases/domain"
)

type OAuthRepository interface {
	GetUserByEmailorUsername(emailOrUsername string) (*domain.Users, error)
	CreateUser(req *domain.Users) error
	GetUserByID(userID uint64) (*domain.Users, error)
	GetUserByIDWithService(userID uint64, serviceID uint64) (*entities.UserWithService, error)
	UpdateUser(req *domain.Users) error
}
