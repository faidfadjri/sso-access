package entities

type Roles struct {
	ServiceRoleId uint64 `gorm:"primaryKey;column:service_role_id" json:"service_role_id"`
	ServiceId     uint64 `gorm:"column:service_id" json:"service_id"`
	RoleName      string `gorm:"column:role_name;size:255" json:"role_name"`
	ServiceName   string `gorm:"column:service_name;size:255" json:"service_name"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at" json:"updated_at"`
}

type AssignedRoles struct {
	UserId        uint64 `gorm:"primaryKey;column:user_id" json:"user_id"`
	ServiceRoleId uint64 `gorm:"primaryKey;column:service_role_id" json:"service_role_id"`
	ServiceId     uint64 `gorm:"primaryKey;column:service_id" json:"service_id"`

	FullName    string `gorm:"column:full_name;size:255" json:"full_name"`
	RoleName    string `gorm:"column:role_name;size:255" json:"role_name"`
	ServiceName string `gorm:"column:service_name;size:255" json:"service_name"`
	CreatedAt   string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   string `gorm:"column:updated_at" json:"updated_at"`
}