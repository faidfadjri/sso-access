package seeding

import (
	"log"

	"akastra-access/internal/infrastructure/databases/domain"

	"gorm.io/gorm"
)

type DepartmentSeeder struct{}

func (s *DepartmentSeeder) Seed(db *gorm.DB) error {
	depts := []string{"IT", "HR", "Finance"}
	for _, d := range depts {
		var count int64
		db.Model(&domain.EmployeeDepartments{}).Where("department_name = ?", d).Count(&count)
		if count == 0 {
			err := db.Create(&domain.EmployeeDepartments{
				DepartmentName: d,
			}).Error
			if err != nil {
				return err
			}
			log.Printf("Seeded Department: %s", d)
		}
	}
	return nil
}
