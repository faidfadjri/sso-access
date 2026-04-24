package domain

type ServiceRoles struct {
	ServiceRoleId uint64 `gorm:"primaryKey;column:service_role_id" json:"service_role_id"`
	ServiceId     uint64 `gorm:"column:service_id" json:"service_id"`
	RoleName      string `gorm:"column:role_name;size:255" json:"role_name"`
	BaseModel
}

func (ServiceRoles) TableName() string {
	return "service_roles"
}

// ----------- ROLE PERMISSION ----------------

type RolePermissions struct {
	RoleId       uint64 `gorm:"primaryKey;column:role_id" json:"role_id"`
	PermissionId uint64 `gorm:"primaryKey;column:permission_id" json:"permission_id"`
	BaseModel
}

func (RolePermissions) TableName() string {
	return "role_permissions"
}

type Permissions struct {
	PermissionId  uint64 `gorm:"primaryKey;column:permission_id" json:"permission_id"`
	PermissionKey string `gorm:"column:permission_key;size:255" json:"permission_key"`
	Description   string `gorm:"column:description;size:255" json:"description"`
	BaseModel
}

func (Permissions) TableName() string {
	return "permissions"
}

// ----------- USER SERVICE ROLE ----------------

type UserServiceRole struct {
	UserId    uint64 `gorm:"primaryKey;column:user_id" json:"user_id"`
	RoleId    uint64 `gorm:"primaryKey;column:role_id" json:"role_id"`
	ServiceId uint64 `gorm:"primaryKey;column:service_id" json:"service_id"`

	BaseModel
}

func (UserServiceRole) TableName() string {
	return "user_service_role"
}

func (UserServiceRole) Relations() []string {
	return []string{"User", "Role", "Service"}
}
