package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateContactMessage(params generated.InsertContactMessageParams, c fiber.Ctx) (generated.ContactMessage, error) {
	return r.Queries.InsertContactMessage(c, params)
}

func (r *Repository) GetContactMessageByID(id uuid.UUID, c fiber.Ctx) (generated.ContactMessage, error) {
	return r.Queries.SelectContactMessageByID(c, id)
}

func (r *Repository) ListContactMessages(params generated.ListContactMessagesParams, c fiber.Ctx) ([]generated.ContactMessage, error) {
	return r.Queries.ListContactMessages(c, params)
}

func (r *Repository) MarkContactMessageAsRead(id uuid.UUID, c fiber.Ctx) (generated.ContactMessage, error) {
	return r.Queries.MarkContactMessageAsRead(c, id)
}

func (r *Repository) DeleteContactMessage(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteContactMessage(c, id)
}
