package request

import "errors"

type ResetPasswordRequest struct {
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (r *ResetPasswordRequest) Validate() error {
	if r.Token == "" {
		return errors.New("token is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.PasswordConfirmation == "" {
		return errors.New("password_confirmation is required")
	}
	if r.Password != r.PasswordConfirmation {
		return errors.New("password and password confirmation do not match")
	}
	return nil
}
