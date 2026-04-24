package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreatePermissionReq struct {
	PermissionKey string `json:"permission_key"`
	Description   string `json:"description"`
}

func (r CreatePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.PermissionKey, validation.Required),
		validation.Field(&r.Description, validation.Required),
	)
}

type UpdatePermissionReq struct {
	PermissionKey string `json:"permission_key"`
	Description   string `json:"description"`
}

func (r UpdatePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.PermissionKey, validation.Required),
		validation.Field(&r.Description, validation.Required),
	)
}
