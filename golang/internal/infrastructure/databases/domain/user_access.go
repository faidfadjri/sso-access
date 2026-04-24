package domain

type UserAccess struct {
	AccessId  uint64 `gorm:"primaryKey;column:access_id"`
	UserId    uint64 `gorm:"column:user_id"`
	ServiceId uint64 `gorm:"column:service_id"`
	Status    string `gorm:"column:status;type:enum('active','revoke')"`

	BaseModel
}

func (UserAccess) TableName() string {
	return "user_access"
}
