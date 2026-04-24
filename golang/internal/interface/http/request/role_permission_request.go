package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateRolePermissionReq struct {
	RoleId       uint64 `json:"role_id"`
	PermissionId uint64 `json:"permission_id"`
}

func (r CreateRolePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RoleId, validation.Required),
		validation.Field(&r.PermissionId, validation.Required),
	)
}

type DeleteRolePermissionReq struct {
	RoleId       uint64 `json:"role_id"`
	PermissionId uint64 `json:"permission_id"`
}

func (r DeleteRolePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RoleId, validation.Required),
		validation.Field(&r.PermissionId, validation.Required),
	)
}
