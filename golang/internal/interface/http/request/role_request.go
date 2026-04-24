package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateRoleReq struct {
	ServiceId uint64 `json:"service_id"`
	RoleName  string `json:"role_name"`
}

func (r CreateRoleReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ServiceId, validation.Required),
		validation.Field(&r.RoleName, validation.Required),
	)
}

type UpdateRoleReq struct {
	ServiceId uint64 `json:"service_id"`
	RoleName  string `json:"role_name"`
}

func (r UpdateRoleReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ServiceId, validation.Required),
		validation.Field(&r.RoleName, validation.Required),
	)
}
