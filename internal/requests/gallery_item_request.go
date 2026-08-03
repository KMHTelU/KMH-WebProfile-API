package requests

import "github.com/google/uuid"

// CreateGalleryItemRequest dipakai alur bulk, yang menyertakan gallery_id di
// dalam body karena tidak ada gallery pada path URL.
type CreateGalleryItemRequest struct {
	GalleryID uuid.UUID `json:"gallery_id" validate:"required"`
	MediaID   uuid.UUID `json:"media_id" validate:"required"`
	SortOrder int32     `json:"sort_order" validate:"omitempty"`
}

type UpdateGalleryItemRequest struct {
	GalleryID uuid.UUID `json:"gallery_id" validate:"required"`
	MediaID   uuid.UUID `json:"media_id" validate:"required"`
	SortOrder int32     `json:"sort_order" validate:"omitempty"`
}
