package requests

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Title           string    `json:"title" validate:"required"`
	Slug            string    `json:"slug" validate:"required"`
	Description     string    `json:"description"`
	EventType       string    `json:"event_type"`
	StartTime       time.Time `json:"start_time" validate:"required"`
	EndTime         time.Time `json:"end_time"`
	Location        string    `json:"location"`
	GoogleMapsUrl   string    `json:"google_maps_url"`
	RegistrationUrl string    `json:"registration_url"`
	CoverMediaID    uuid.UUID `json:"cover_media_id"`
	Status          string    `json:"status"`
	IsPublished     bool      `json:"is_published"`
}

type UpdateEventRequest struct {
	Title           string    `json:"title" validate:"required"`
	Slug            string    `json:"slug" validate:"required"`
	Description     string    `json:"description"`
	EventType       string    `json:"event_type"`
	StartTime       time.Time `json:"start_time" validate:"required"`
	EndTime         time.Time `json:"end_time"`
	Location        string    `json:"location"`
	GoogleMapsUrl   string    `json:"google_maps_url"`
	RegistrationUrl string    `json:"registration_url"`
	CoverMediaID    uuid.UUID `json:"cover_media_id"`
	Status          string    `json:"status"`
	IsPublished     bool      `json:"is_published"`
}
