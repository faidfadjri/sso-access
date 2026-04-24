package request

import "errors"

type ForgotPasswordRequest struct {
	Email      string `json:"email"`
	ForgotType string `json:"forgot_type"`
}

func (r *ForgotPasswordRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.ForgotType != "username" && r.ForgotType != "password" {
		return errors.New("forgot_type must be username or password")
	}
	return nil
}
