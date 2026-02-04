package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
)

func (r *Repository) InsertLog(params generated.InsertActivityLogParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertActivityLog(c, params)
	if err != nil {
		return err
	}
	return nil
}
