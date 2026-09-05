package requests

import "github.com/go-playground/validator/v10"

type FindAllRoles struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	ID   *int   `json:"id"`
	Name string `json:"name" validate:"required"`
}

// RoleRequestPayload is a Kafka payload for role validation requests.
type RoleRequestPayload struct {
	UserID        int    `json:"user_id"`
	CorrelationID string `json:"correlation_id"`
	ReplyTopic    string `json:"reply_topic"`
}

func (r *CreateRoleRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

func (r *UpdateRoleRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}
