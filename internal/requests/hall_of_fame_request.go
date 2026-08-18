package requests

import "github.com/google/uuid"

// DTO Hall of Fame. Semua field teks opsional selain yang divalidasi required;
// slice kosong berarti daftar sengaja dikosongkan.

type HofGenerationRequest struct {
	Name        string   `json:"name" validate:"required"`
	YearStart   int32    `json:"year_start" validate:"required"`
	YearEnd     int32    `json:"year_end" validate:"required"`
	Description string   `json:"description"`
	Milestones  []string `json:"milestones"`
	Accent      string   `json:"accent"`
	SortOrder   int32    `json:"sort_order"`
}

type HofPersonRequest struct {
	GenerationID  uuid.UUID `json:"generation_id" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Role          string    `json:"role"`
	StudyProgram  string    `json:"study_program"`
	Biography     string    `json:"biography"`
	Contributions string    `json:"contributions"`
	Legacy        string    `json:"legacy"`
	// Kutipan HARUS terverifikasi; kosongkan bila tidak ada.
	Quote        string    `json:"quote"`
	Fields       []string  `json:"fields"`
	PhotoMediaID uuid.UUID `json:"photo_media_id"`
	SortOrder    int32     `json:"sort_order"`
}

type HofAchievementRequest struct {
	PersonID     uuid.UUID `json:"person_id" validate:"required"`
	Title        string    `json:"title" validate:"required"`
	Category     string    `json:"category"`
	Year         int32     `json:"year" validate:"required"`
	Organization string    `json:"organization"`
	Result       string    `json:"result"`
	Description  string    `json:"description"`
}

type HofTimelineEventRequest struct {
	YearLabel   string `json:"year_label" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sort_order"`
}
