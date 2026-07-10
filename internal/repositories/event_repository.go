package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateEvent(params generated.InsertEventParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertEvent(c, params)
	return err
}

func (r *Repository) GetEventByID(id uuid.UUID, c fiber.Ctx) (generated.GetEventByIDRow, error) {
	return r.Queries.GetEventByID(c, id)
}

func (r *Repository) ListEvents(params generated.ListEventsParams, c fiber.Ctx) ([]generated.ListEventsRow, error) {
	return r.Queries.ListEvents(c, params)
}

func (r *Repository) UpdateEvent(params generated.UpdateEventParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateEvent(c, params)
	return err
}

func (r *Repository) DeleteEvent(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteEvent(c, id)
}
