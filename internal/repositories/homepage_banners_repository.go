package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) InsertHomepageBanner(params generated.InsertHomepageBannerParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertHomepageBanner(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetHomepageBanners(params generated.SelectAllHomepageBannersParams, c fiber.Ctx) ([]generated.SelectAllHomepageBannersRow, error) {
	rows, err := r.Queries.SelectAllHomepageBanners(c, params)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) DeleteHomepageBanner(id uuid.UUID, c fiber.Ctx) error {
	err := r.Queries.DeleteHomepageBanner(c, id)
	if err != nil {
		return err
	}
	return nil
}
