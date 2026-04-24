package request

import "errors"

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthorizeRequest struct {
	ClientID     		string `json:"client_id"`
	RedirectURI  		string `json:"redirect_uri"`
	ResponseType 		string `json:"response_type"`
	Scope        		string `json:"scope"`
	State        		string `json:"state"`
	CodeChallenge 		*string `json:"code_challenge"`
	CodeChallengeMethod *string `json:"code_challenge_method"`
}

func (r *AuthorizeRequest) Validate() error {
	if r.ClientID == "" {
		return errors.New("client_id is required")
	}

	if r.RedirectURI == "" {
		return errors.New("redirect_uri is required")
	}

	if r.ResponseType == "" {
		return errors.New("response_type is required")
	}

	if r.Scope == "" {
		return errors.New("scope is required")
	}

	if r.State == "" {
		return errors.New("state is required")
	}

	return nil
}

type ExchangeAuthCodeRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret *string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	CodeVerifier string `json:"code_verifier"` // For PKCE support later
}

func (r *ExchangeAuthCodeRequest) Validate() error {
	if r.GrantType != "authorization_code" {
		return errors.New("invalid grant_type")
	}
	if r.Code == "" {
		return errors.New("code is required")
	}
	if r.ClientID == "" {
		return errors.New("client_id is required")
	}
	if r.RedirectURI == "" {
		return errors.New("redirect_uri is required")
	}
	return nil
}

type UpdateAccountRequest struct {
	FullName 				string `json:"full_name"`
	Email 					string `json:"email"`
	Username 				string `json:"username"`
	Phone 					string `json:"phone"`
	Password 				*string `json:"password"`
	PasswordConfirmation 	*string `json:"password_confirmation"`
}

func (r *UpdateAccountRequest) Validate() error {
	if r.FullName == "" {
		return errors.New("full_name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Phone == "" {
		return errors.New("phone is required")
	}
	if r.Password != nil && *r.Password != "" {
		if r.PasswordConfirmation == nil || *r.PasswordConfirmation == "" {
			return errors.New("password_confirmation is required")
		}
		if *r.Password != *r.PasswordConfirmation {
			return errors.New("password and password_confirmation do not match")
		}
	}

	return nil
}