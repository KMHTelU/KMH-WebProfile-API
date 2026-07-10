package requests

type CreateOrganizationProfileRequest struct {
	Name         string `json:"name" validate:"required"`
	ShortName    string `json:"short_name" validate:"required"`
	Description  string `json:"description"`
	Vision       string `json:"vision"`
	Mission      string `json:"mission"`
	History      string `json:"history"`
	Address      string `json:"address"`
	Email        string `json:"email" validate:"omitempty,email"`
	Phone        string `json:"phone"`
	InstagramUrl string `json:"instagram_url" validate:"omitempty,url"`
	YoutubeUrl   string `json:"youtube_url" validate:"omitempty,url"`
	WebsiteUrl   string `json:"website_url" validate:"omitempty,url"`
}

type UpdateOrganizationProfileRequest struct {
	Name         string `json:"name" validate:"required"`
	ShortName    string `json:"short_name" validate:"required"`
	Description  string `json:"description"`
	Vision       string `json:"vision"`
	Mission      string `json:"mission"`
	History      string `json:"history"`
	Address      string `json:"address"`
	Email        string `json:"email" validate:"omitempty,email"`
	Phone        string `json:"phone"`
	InstagramUrl string `json:"instagram_url" validate:"omitempty,url"`
	YoutubeUrl   string `json:"youtube_url" validate:"omitempty,url"`
	WebsiteUrl   string `json:"website_url" validate:"omitempty,url"`
}
