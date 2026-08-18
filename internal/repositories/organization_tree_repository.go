package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
)

func (r *Repository) ListOrgTreeDivisions(c fiber.Ctx) ([]generated.SelectOrgTreeDivisionsRow, error) {
	return r.Queries.SelectOrgTreeDivisions(c)
}

func (r *Repository) ListOrgTreeAssignments(c fiber.Ctx) ([]generated.SelectOrgTreeAssignmentsRow, error) {
	return r.Queries.SelectOrgTreeAssignments(c)
}
