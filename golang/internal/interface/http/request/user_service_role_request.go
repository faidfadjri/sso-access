package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateUserServiceRoleReq struct {
	RoleId    uint64   `json:"role_id"`
	UserIds   []uint64 `json:"user_ids"`
}

func (r CreateUserServiceRoleReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserIds, validation.Required),
		validation.Field(&r.RoleId, validation.Required),
	)
}

type DeleteUserServiceRoleReq struct {
	UserId    uint64 `json:"user_id"`
	ServiceId uint64 `json:"service_id"`
	RoleId    uint64 `json:"role_id"`
}

func (r DeleteUserServiceRoleReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserId, validation.Required),
		validation.Field(&r.ServiceId, validation.Required),
		validation.Field(&r.RoleId, validation.Required),
	)
}
