package domain

type Services struct {
	ServiceId    uint64  `gorm:"primaryKey;column:service_id" json:"service_id"`
	ServiceName  string  `gorm:"column:service_name;size:255" json:"service_name"`
	Description  *string `gorm:"column:description;size:255" json:"description"`
	Logo         *string `gorm:"column:logo;size:255" json:"logo"`
	ClientId     string  `gorm:"column:client_id;size:255" json:"client_id"`
	ClientSecret string  `gorm:"column:client_secret;size:255" json:"client_secret"`
	RedirectUrl  string  `gorm:"column:redirect_url;size:255" json:"redirect_url"`
	IsActive     bool    `gorm:"column:is_active" json:"is_active"`
	BaseModel
}

func (Services) TableName() string {
	return "services"
}
