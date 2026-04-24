package seeding

import (
	"log"

	"akastra-access/internal/infrastructure/databases/domain"

	"gorm.io/gorm"
)

type PositionSeeder struct{}

func (s *PositionSeeder) Seed(db *gorm.DB) error {
	positions := []string{"Manager", "Staff", "Developer"}
	for _, p := range positions {
		var count int64
		db.Model(&domain.EmployeePositions{}).Where("position_name = ?", p).Count(&count)
		if count == 0 {
			err := db.Create(&domain.EmployeePositions{
				PositionName: p,
			}).Error
			if err != nil {
				return err
			}
			log.Printf("Seeded Position: %s", p)
		}
	}
	return nil
}
