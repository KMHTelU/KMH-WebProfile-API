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

// Nama field form sengaja datar (tanpa prefix "data.") karena decoder form
// Fiber memperlakukan titik sebagai path struct bersarang, sehingga key
// seperti "data.title" tidak pernah cocok dan seluruh field dibiarkan kosong.
// IsActive tidak boleh divalidasi "required" karena nilai false dianggap kosong
// oleh validator dan membuat banner nonaktif mustahil dibuat.
type HomepageBannerRequest struct {
	Title     string    `form:"title" json:"title" validate:"required"`
	Subtitle  string    `form:"subtitle" json:"subtitle" validate:"omitempty"`
	CtaText   string    `form:"cta_text" json:"cta_text" validate:"omitempty"`
	CtaUrl    string    `form:"cta_url" json:"cta_url" validate:"omitempty,url"`
	IsActive  bool      `form:"is_active" json:"is_active"`
	StartDate time.Time `form:"start_date" json:"start_date" validate:"required"`
	EndDate   time.Time `form:"end_date" json:"end_date" validate:"required"`
	AltText   string    `form:"alt_text" json:"alt_text" validate:"omitempty"`
	Caption   string    `form:"caption" json:"caption" validate:"omitempty"`
}
