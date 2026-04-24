package seeding

import (
	"log"

	"akastra-access/internal/infrastructure/databases/domain"

	"gorm.io/gorm"
)

type RoleSeeder struct{}

func (s *RoleSeeder) Seed(db *gorm.DB) error {
	roles := []string{"Super Admin", "User", "Employee"}
	
	// Ensure we have a service to attach roles to
	var service domain.Services
	if err := db.First(&service).Error; err != nil {
		log.Println("Skipping role seeding: No services found")
		return nil
	}

	for _, roleName := range roles {
		var count int64
		// Check if role exists for this service
		db.Model(&domain.ServiceRoles{}).Where("role_name = ? AND service_id = ?", roleName, service.ServiceId).Count(&count)
		if count == 0 {
			err := db.Create(&domain.ServiceRoles{
				ServiceId: service.ServiceId,
				RoleName:  roleName,
			}).Error
			if err != nil {
				return err
			}
			log.Printf("Seeded Role: %s", roleName)
		}
	}
	return nil
}
