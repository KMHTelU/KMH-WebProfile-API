package requests

import "github.com/google/uuid"

type CreateMemberDivisionRequest struct {
	MemberID   uuid.UUID `json:"member_id" validate:"required"`
	DivisionID uuid.UUID `json:"division_id" validate:"required"`
	RoleTitle  string    `json:"role_title" validate:"omitempty"`
}

type UpdateMemberDivisionRequest struct {
	MemberID   uuid.UUID `json:"member_id" validate:"required"`
	DivisionID uuid.UUID `json:"division_id" validate:"required"`
	RoleTitle  string    `json:"role_title" validate:"omitempty"`
}
