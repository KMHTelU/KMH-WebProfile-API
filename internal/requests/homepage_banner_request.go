package requests

import "time"

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
