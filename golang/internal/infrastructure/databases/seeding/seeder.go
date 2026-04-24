package seeding

import (
	"log"

	"gorm.io/gorm"
)

type Seeder interface {
	Seed(db *gorm.DB) error
}

func Run(db *gorm.DB) {
	log.Println("Seeding database...")

	seeders := []Seeder{
		&ServiceSeeder{},
		&RoleSeeder{},
		&UserSeeder{},
	}

	for _, seeder := range seeders {
		if err := seeder.Seed(db); err != nil {
			log.Printf("Failed to seed: %v", err)
		}
	}

	log.Println("Database seeding completed successfully.")
}
