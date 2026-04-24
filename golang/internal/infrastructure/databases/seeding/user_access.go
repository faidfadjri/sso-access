package seeding

import (
	"log"

	"akastra-access/internal/infrastructure/databases/domain"

	"gorm.io/gorm"
)

type UserAccessSeeder struct{}

func (s *UserAccessSeeder) Seed(db *gorm.DB) error {
	var count int64
	db.Model(&domain.UserAccess{}).Count(&count)
	if count > 0 {
		return nil
	}

	userAccess := domain.UserAccess{
		UserId:    1,
		ServiceId: 1,
		Status:    "active",
	}

	if err := db.Create(&userAccess).Error; err != nil {
		return err
	}

	log.Println("Seeded User Access")
	return nil
}
