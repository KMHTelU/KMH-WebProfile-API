package requests

import "github.com/google/uuid"

type CreateGalleryRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	EventID     uuid.UUID `json:"event_id"`
	IsPublic    bool      `json:"is_public"`
}

type UpdateGalleryRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	EventID     uuid.UUID `json:"event_id"`
	IsPublic    bool      `json:"is_public"`
}

type AddGalleryItemRequest struct {
	MediaID   uuid.UUID `json:"media_id" validate:"required"`
	SortOrder int32     `json:"sort_order"`
}
