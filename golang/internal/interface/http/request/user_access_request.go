package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateUserAccessReq struct {
	UserId     uint64   `json:"user_id"`
	ServiceIds []uint64 `json:"service_ids"`
	Status     string   `json:"status"`
}

func (r CreateUserAccessReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserId, validation.Required),
		validation.Field(&r.ServiceIds, validation.Required),
		validation.Field(&r.Status, validation.Required, validation.In("active", "revoke")),
	)
}

type UpdateUserAccessReq struct {
	UserId     uint64   `json:"user_id"`
	ServiceIds []uint64 `json:"service_ids"`
	Status     string   `json:"status"`
}

func (r UpdateUserAccessReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserId, validation.Required),
		validation.Field(&r.ServiceIds, validation.Required),
		validation.Field(&r.Status, validation.Required, validation.In("active", "revoke")),
	)
}
