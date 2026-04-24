package seeding

import (
	"log"

	"akastra-access/internal/infrastructure/databases/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (s *UserSeeder) Seed(db *gorm.DB) error {
	var count int64
	email := "email@gmail.com"
	db.Model(&domain.Users{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("very-secret"), bcrypt.DefaultCost)
	
	admin := domain.Users{
		FullName:  "Super Administrator",
		Email:     email,
		Username:  "super-admin",
		Password:  string(hashedPassword),
		Admin:     true,
	}
	
	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to seed admin user: %v", err)
		return err
	}
	log.Println("Seeded Admin User")

	s.seedEmployee(db, admin.UserId)
	s.seedUserRole(db, admin.UserId)
	return nil
}

func (s *UserSeeder) seedEmployee(db *gorm.DB, userId uint64) {
	var dept domain.EmployeeDepartments
	db.Where("department_name = ?", "IT").First(&dept)
	
	var pos domain.EmployeePositions
	db.Where("position_name = ?", "Developer").First(&pos)

	if dept.DepartmentId != 0 && pos.PositionId != 0 {
		emp := domain.Employees{
			UserId:       userId,
			DepartmentId: dept.DepartmentId,
			PositionId:   pos.PositionId,
			Nrp:          1001,
		}
		if err := db.Create(&emp).Error; err != nil {
			log.Printf("Failed to seed employee data: %v", err)
		} else {
			log.Println("Seeded Admin Employee Data")
		}
	}
}

func (s *UserSeeder) seedUserRole(db *gorm.DB, userId uint64) {
	var service domain.Services
	db.First(&service)
	
	var role domain.ServiceRoles
	db.Where("role_name = ? AND service_id = ?", "Super Admin", service.ServiceId).First(&role)

	if service.ServiceId != 0 && role.ServiceRoleId != 0 {
		usr := domain.UserServiceRole{
			UserId:    userId,
			RoleId:    role.ServiceRoleId,
			ServiceId: service.ServiceId,
		}
		if err := db.Create(&usr).Error; err != nil {
			log.Printf("Failed to seed user role: %v", err)
		} else {
			log.Println("Seeded User Service Role")
		}
	}
}
