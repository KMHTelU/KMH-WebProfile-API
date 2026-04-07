package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateDivision(params generated.InsertDivisionParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertDivision(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetDivisionByID(id uuid.UUID, c fiber.Ctx) (generated.GetDivisionByIDRow, error) {
	division, err := r.Queries.GetDivisionByID(c, id)
	if err != nil {
		return generated.GetDivisionByIDRow{}, err
	}
	return division, nil
}

func (r *Repository) GetAllDivisions(c fiber.Ctx) ([]generated.GetAllDivisionsRow, error) {
	divisions, err := r.Queries.GetAllDivisions(c)
	if err != nil {
		return nil, err
	}
	return divisions, nil
}

func (r *Repository) UpdateDivision(params generated.UpdateDivisionParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateDivision(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteDivision(id uuid.UUID, c fiber.Ctx) error {
	err := r.Queries.DeleteDivision(c, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateDivisionIcon(params generated.UpdateDivisionIconParams, c fiber.Ctx) error {
	err := r.Queries.UpdateDivisionIcon(c, params)
	if err != nil {
		return err
	}
	return nil
}
