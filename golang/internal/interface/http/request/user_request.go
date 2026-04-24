package request

import "fmt"

type CreateUserReq struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone"`
	Admin    string `json:"admin"` // "true" or "false" from form-data
}

func (r *CreateUserReq) Validate() error {
	if r.FullName == "" {
		return fmt.Errorf("full_name is required")
	}
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if r.Username == "" {
		return fmt.Errorf("username is required")
	}
	if r.Password == "" || len(r.Password) < 6 {
		return fmt.Errorf("password is required and must be at least 6 characters")
	}
	return nil
}

type UpdateUserReq struct {
	ID              uint64 `json:"id"`
	FullName        string `json:"full_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Phone           string `json:"phone"`
	Admin           string `json:"admin"`
}

func (r *UpdateUserReq) Validate() error {
	if r.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if r.FullName == "" {
		return fmt.Errorf("full_name is required")
	}
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if r.Username == "" {
		return fmt.Errorf("username is required")
	}
	if r.Password != "" {
		if len(r.Password) < 6 {
			return fmt.Errorf("password must be at least 6 characters")
		}
		if r.Password != r.PasswordConfirm {
			return fmt.Errorf("password confirmation does not match")
		}
	}
	return nil
}

type BatchDeleteUserReq struct {
	IDs []uint64 `json:"ids" validate:"required"`
}
