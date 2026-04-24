package seeding

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"akastra-access/internal/app/config"
	"akastra-access/internal/infrastructure/databases/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ServiceSeeder struct{}

func (s *ServiceSeeder) Seed(db *gorm.DB) error {
	var count int64
	db.Model(&domain.Services{}).Count(&count)
	if count > 0 {
		return nil
	}

	clientID, err := generateClientID()
	if err != nil {
		return err
	}

	clientSecret, err := generateClientSecret()
	if err != nil {
		return err
	}

	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	service := domain.Services{
		ServiceName:  "SSO System",
		ClientId:     clientID,
		ClientSecret: string(hashedSecret),
		RedirectUrl:  config.Load().FrontendURL,
		IsActive:     true,
	}
	if err := db.Create(&service).Error; err != nil {
		return err
	}
	log.Println("Seeded Services")
	return nil
}


func generateClientID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateClientSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
