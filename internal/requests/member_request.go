package requests

type CreateMemberRequest struct {
	Name         string `json:"name" validate:"required"`
	Nim          string `json:"nim" validate:"required"`
	Bio          string `json:"bio" validate:"omitempty"`
	Email        string `json:"email" validate:"omitempty,email"`
	Phone        string `json:"phone" validate:"omitempty"`
	InstagramUrl string `json:"instagram_url" validate:"omitempty,url"`
	Faculty      string `json:"faculty" validate:"omitempty"`
	StudyProgram string `json:"study_program" validate:"omitempty"`
	CohortYear   int32  `json:"cohort_year" validate:"omitempty"`
	PeriodStart  int32  `json:"period_start" validate:"required"`
	PeriodEnd    int32  `json:"period_end" validate:"required"`
}

type UpdateMemberRequest struct {
	Name         string `json:"name" validate:"omitempty"`
	Nim          string `json:"nim" validate:"omitempty"`
	Bio          string `json:"bio" validate:"omitempty"`
	Email        string `json:"email" validate:"omitempty,email"`
	Phone        string `json:"phone" validate:"omitempty"`
	InstagramUrl string `json:"instagram_url" validate:"omitempty,url"`
	Faculty      string `json:"faculty" validate:"omitempty"`
	StudyProgram string `json:"study_program" validate:"omitempty"`
	CohortYear   int32  `json:"cohort_year" validate:"omitempty"`
	PeriodStart  int32  `json:"period_start" validate:"omitempty"`
	PeriodEnd    int32  `json:"period_end" validate:"omitempty"`
	IsActive     *bool  `json:"is_active" validate:"omitempty"`
	PhotoMediaID string `json:"photo_media_id" validate:"omitempty,uuid4"`
}
