package request

import "fmt"

type CreateClientReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Logo        string `json:"logo"` // Will hold the path after upload
	RedirectURL string `json:"redirect_url" validate:"required"`
}

func (r *CreateClientReq) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	// Logo is handled separately
	if r.RedirectURL == "" {
		return fmt.Errorf("redirect_url is required")
	}
	return nil
}

type UpdateClientReq struct {
	ID          uint64 `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Logo        string `json:"logo"`
	RedirectURL string `json:"redirect_url" validate:"required"`
	IsActive    *bool  `json:"is_active"`
}

func (r *UpdateClientReq) Validate() error {
	if r.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	if r.RedirectURL == "" {
		return fmt.Errorf("redirect_url is required")
	}
	return nil
}