package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateRole(params generated.InsertRoleParams, c fiber.Ctx) (generated.Role, error) {
	return r.Queries.InsertRole(c, params)
}

func (r *Repository) GetRoleByID(id uuid.UUID, c fiber.Ctx) (generated.Role, error) {
	return r.Queries.GetRoleByID(c, id)
}

func (r *Repository) GetAllRoles(c fiber.Ctx) ([]generated.Role, error) {
	return r.Queries.GetAllRoles(c)
}

func (r *Repository) UpdateRole(params generated.UpdateRoleParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateRole(c, params)
	return err
}

func (r *Repository) DeleteRole(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteRole(c, id)
}
