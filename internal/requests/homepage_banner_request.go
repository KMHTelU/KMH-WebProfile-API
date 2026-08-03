package requests

import (
	"time"

	"github.com/google/uuid"
)

// HomepageBannerJSONRequest dipakai alur bulk. Berbeda dari HomepageBannerRequest
// yang mengunggah berkas lewat multipart, di sini gambar dirujuk lewat media_id
// hasil unggahan sebelumnya dan boleh dikosongkan.
type HomepageBannerJSONRequest struct {
	Title     string    `json:"title" validate:"required"`
	Subtitle  string    `json:"subtitle" validate:"omitempty"`
	MediaID   uuid.UUID `json:"media_id" validate:"omitempty"`
	CtaText   string    `json:"cta_text" validate:"omitempty"`
	CtaUrl    string    `json:"cta_url" validate:"omitempty,url"`
	IsActive  bool      `json:"is_active" validate:"omitempty"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

type HomepageBannerRequest struct {
	Title     string    `form:"data.title" json:"title" validate:"required"`
	Subtitle  string    `form:"data.subtitle" json:"subtitle" validate:"omitempty"`
	CtaText   string    `form:"data.cta_text" json:"cta_text" validate:"omitempty"`
	CtaUrl    string    `form:"data.cta_url" json:"cta_url" validate:"omitempty,url"`
	IsActive  bool      `form:"data.is_active" json:"is_active" validate:"required"`
	StartDate time.Time `form:"data.start_date" json:"start_date" validate:"required"`
	EndDate   time.Time `form:"data.end_date" json:"end_date" validate:"required"`
	AltText   string    `form:"data.alt_text" json:"alt_text" validate:"omitempty"`
	Caption   string    `form:"data.caption" json:"caption" validate:"omitempty"`
}
