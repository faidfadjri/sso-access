package domain

type Users struct {
	UserId   uint64  `gorm:"primaryKey;column:user_id" json:"user_id"`
	FullName string  `gorm:"column:full_name;size:255" json:"full_name"`
	Email    string  `gorm:"column:email;size:255" json:"email"`
	Username string  `gorm:"column:username;size:255" json:"username"`
	Password string  `gorm:"column:password;size:255" json:"password"`
	Photo    *string `gorm:"column:photo;size:255" json:"photo"`
	Phone    *string `gorm:"column:phone;size:255" json:"phone"`
	Admin    bool    `gorm:"column:admin;default:false" json:"admin"`
	BaseModel
}

func (Users) TableName() string {
	return "users"
}

func (u *Users) GetUserID() uint64       { return u.UserId }
func (u *Users) GetEmail() string        { return u.Email }
func (u *Users) GetUsername() string     { return u.Username }
func (u *Users) GetFullName() string     { return u.FullName }
func (u *Users) GetPhone() *string       { return u.Phone }
func (u *Users) GetPhoto() *string       { return u.Photo }
func (u *Users) GetServiceName() *string { return nil }
func (u *Users) GetRoleName() *string    { return nil }
func (u *Users) GetIDPRole() *string     { return nil }
