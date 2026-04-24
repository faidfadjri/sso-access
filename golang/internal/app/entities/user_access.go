package entities

type UserAccess struct {
	AccessId  uint64 `json:"access_id"`
	UserId    uint64 `json:"user_id"`
	ServiceId uint64 `json:"service_id"`
	Status    string `json:"status"`

	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`

	ServiceName string `json:"service_name"`
	RedirectURL string `json:"redirect_url"`
	ClientId    string `json:"client_id"`
	CreatedAt   string `json:"created_at"`
}