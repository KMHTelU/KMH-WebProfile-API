package requests

import "github.com/google/uuid"

type CreateDivisionRequest struct {
	Name          string    `json:"name" validate:"required"`
	Slug          string    `json:"slug" validate:"required"`
	Description   string    `json:"description" validate:"omitempty"`
	CoordinatorID uuid.UUID `json:"coordinator_id" validate:"omitempty,uuid4"`
}

type UpdateDivisionRequest struct {
	Name          string    `json:"name" validate:"omitempty"`
	Slug          string    `json:"slug" validate:"omitempty"`
	Description   string    `json:"description" validate:"omitempty"`
	CoordinatorID uuid.UUID `json:"coordinator_id" validate:"omitempty,uuid4"`
	IsActive      bool      `json:"is_active" validate:"omitempty"`
}
