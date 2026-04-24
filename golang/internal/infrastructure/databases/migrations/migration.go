package migrations

import (
	"log"

	"akastra-access/internal/infrastructure/databases/domain"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// AutoMigrate all models
	err := db.AutoMigrate(
		&domain.EmployeeDepartments{},  
		&domain.EmployeePositions{},
		&domain.Employees{},
		&domain.Users{},
		&domain.Services{},
		&domain.Permissions{},
		&domain.ServiceRoles{},
		&domain.RolePermissions{},
		&domain.UserServiceRole{},
		&domain.UserAccess{},
		&domain.ActivityLogs{},       
	)


	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	// Manually create foreign key constraints for tables (except employees, departments, positions and logs)
	constraints := []string{
		"ALTER TABLE service_roles ADD CONSTRAINT fk_service_roles_service FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE CASCADE ON UPDATE CASCADE",
		"ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES service_roles(service_role_id) ON DELETE CASCADE ON UPDATE CASCADE",
		"ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE ON UPDATE CASCADE",
		"ALTER TABLE user_service_role ADD CONSTRAINT fk_user_service_role_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE",
		"ALTER TABLE user_service_role ADD CONSTRAINT fk_user_service_role_role FOREIGN KEY (role_id) REFERENCES service_roles(service_role_id) ON DELETE CASCADE ON UPDATE CASCADE",
		"ALTER TABLE user_service_role ADD CONSTRAINT fk_user_service_role_service FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE CASCADE ON UPDATE CASCADE",
		"ALTER TABLE user_access ADD CONSTRAINT fk_user_access_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE",
		"ALTER TABLE user_access ADD CONSTRAINT fk_user_access_service FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE CASCADE ON UPDATE CASCADE",
	}

	for _, query := range constraints {
		if err := db.Exec(query).Error; err != nil {
			log.Printf("Warning setting constraint: %v", err)
		}
	}

	log.Println("Migration successful!")
	return nil
}
