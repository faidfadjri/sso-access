package entities

import "time"

type UserWithService struct {
	UserId      uint64     `gorm:"primaryKey;column:user_id" json:"user_id"`
	FullName    string     `gorm:"column:full_name;size:255" json:"full_name"`
	Email       string     `gorm:"column:email;size:255" json:"email"`
	Username    string     `gorm:"column:username;size:255" json:"username"`
	Photo       *string    `gorm:"column:photo;size:255" json:"photo"`
	Phone       *string    `gorm:"column:phone;size:255" json:"phone"`
	ServiceName *string    `gorm:"column:service_name;size:255" json:"service_name"`
	RoleName    *string    `gorm:"column:role_name;size:255" json:"role_name"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (u *UserWithService) GetUserID() uint64       { return u.UserId }
func (u *UserWithService) GetEmail() string        { return u.Email }
func (u *UserWithService) GetUsername() string     { return u.Username }
func (u *UserWithService) GetFullName() string     { return u.FullName }
func (u *UserWithService) GetPhone() *string       { return u.Phone }
func (u *UserWithService) GetPhoto() *string       { return u.Photo }
func (u *UserWithService) GetServiceName() *string { return u.ServiceName }
func (u *UserWithService) GetRoleName() *string    { return u.RoleName }
func (u *UserWithService) GetIDPRole() *string     { return nil }