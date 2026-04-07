package requests

type CreateMemberRequest struct {
	Name         string `json:"name" validate:"required"`
	Npm          string `json:"npm" validate:"required"`
	Bio          string `json:"bio" validate:"omitempty"`
	Email        string `json:"email" validate:"omitempty,email"`
	Phone        string `json:"phone" validate:"omitempty"`
	InstagramUrl string `json:"instagram_url" validate:"omitempty,url"`
	PeriodStart  int32  `json:"period_start" validate:"required"`
	PeriodEnd    int32  `json:"period_end" validate:"required"`
}

type UpdateMemberRequest struct {
	Name         string `json:"name" validate:"omitempty"`
	Npm          string `json:"npm" validate:"omitempty"`
	Bio          string `json:"bio" validate:"omitempty"`
	Email        string `json:"email" validate:"omitempty,email"`
	Phone        string `json:"phone" validate:"omitempty"`
	InstagramUrl string `json:"instagram_url" validate:"omitempty,url"`
	PeriodStart  int32  `json:"period_start" validate:"omitempty"`
	PeriodEnd    int32  `json:"period_end" validate:"omitempty"`
	IsActive     bool   `json:"is_active" validate:"omitempty"`
	PhotoMediaID string `json:"photo_media_id" validate:"omitempty,uuid4"`
}
