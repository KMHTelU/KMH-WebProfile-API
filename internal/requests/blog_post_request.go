package requests

import "github.com/google/uuid"

type CreateBlogPostRequest struct {
	Title           string      `json:"title" validate:"required"`
	Slug            string      `json:"slug" validate:"required"`
	Excerpt         string      `json:"excerpt"`
	Content         string      `json:"content" validate:"required"`
	CategoryID      uuid.UUID   `json:"category_id"`
	FeaturedMediaID uuid.UUID   `json:"featured_media_id"`
	Status          string      `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	TagIDs          []uuid.UUID `json:"tag_ids"`
}

type UpdateBlogPostRequest struct {
	Title           string      `json:"title" validate:"required"`
	Slug            string      `json:"slug" validate:"required"`
	Excerpt         string      `json:"excerpt"`
	Content         string      `json:"content" validate:"required"`
	CategoryID      uuid.UUID   `json:"category_id"`
	FeaturedMediaID uuid.UUID   `json:"featured_media_id"`
	Status          string      `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	TagIDs          []uuid.UUID `json:"tag_ids"`
}
