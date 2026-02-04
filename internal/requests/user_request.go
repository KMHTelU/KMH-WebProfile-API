package requests

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8"`
	RoleID   uuid.UUID `json:"role_id" validate:"required"`
}

type UpdateUserRequest struct {
	Name     string    `json:"name" validate:"omitempty"`
	Email    string    `json:"email" validate:"omitempty,email"`
	Password string    `json:"password" validate:"omitempty,min=8"`
	IsActive bool      `json:"is_active" validate:"omitempty"`
	RoleID   uuid.UUID `json:"role_id" validate:"omitempty"`
}
