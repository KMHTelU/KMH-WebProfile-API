package entities

import "github.com/google/uuid"

type HomepageBanner struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	CtaText   string    `json:"cta_text"`
	CtaUrl    string    `json:"cta_url"`
	IsActive  bool      `json:"is_active"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	Media     *Media    `json:"media,omitempty"`
}
