package requests

type CreateContactMessageRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Subject string `json:"subject"`
	Message string `json:"message" validate:"required"`
}
