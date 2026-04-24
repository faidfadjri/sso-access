package response

type LoginResponse struct {
	SessionID string `json:"session_id"`
	Fullname  string `json:"full_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Photo     string `json:"photo"`
	Phone     string `json:"phone"`
	IsAdmin   bool   `json:"is_admin"`
}

type LoginResponseJWT struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// MeResponse is the payload returned by the /me endpoint.
// It is sourced from the JWT AccessClaims so it includes service and role context.
type MeResponse struct {
	UserID      uint64  `json:"user_id"`
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	Phone       *string `json:"phone"`
	Photo       *string `json:"photo"`
	ServiceName *string `json:"service_name"`
	RoleName    *string `json:"role_name"`
}