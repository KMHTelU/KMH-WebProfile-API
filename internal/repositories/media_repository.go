package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) InsertMedia(params generated.InsertMediaParams, c fiber.Ctx) (generated.Medium, error) {
	media, err := r.Queries.InsertMedia(c, params)
	if err != nil {
		return generated.Medium{}, err
	}
	return media, nil
}

func (r *Repository) DeleteMedia(id uuid.UUID, c fiber.Ctx) error {
	err := r.Queries.DeleteMedia(c, id)
	if err != nil {
		return err
	}
	return nil
}
