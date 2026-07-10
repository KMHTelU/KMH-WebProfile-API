package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateOrganizationProfile(params generated.InsertOrganizationProfileParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertOrganizationProfile(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateOrganizationProfile(params generated.UpdateOrganizationProfileParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateOrganizationProfile(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateOrganizationProfileLogo(params generated.UpdateOrganizationProfileLogoParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateOrganizationProfileLogo(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetOrganizationProfile(id uuid.UUID, c fiber.Ctx) (generated.GetOrganizationProfileRow, error) {
	profile, err := r.Queries.GetOrganizationProfile(c, id)
	if err != nil {
		return generated.GetOrganizationProfileRow{}, err
	}
	return profile, nil
}

func (r *Repository) DeleteOrganizationProfile(id uuid.UUID, c fiber.Ctx) error {
	err := r.Queries.DeleteOrganizationProfile(c, id)
	if err != nil {
		return err
	}
	return nil
}
