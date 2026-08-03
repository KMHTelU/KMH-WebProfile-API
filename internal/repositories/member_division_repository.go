package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateMemberDivision(params generated.InsertMemberDivisionParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertMemberDivision(c, params)
	return err
}

func (r *Repository) UpdateMemberDivision(params generated.UpdateMemberDivisionParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateMemberDivision(c, params)
	return err
}

func (r *Repository) GetMemberDivisionByID(id uuid.UUID, c fiber.Ctx) (generated.MemberDivision, error) {
	return r.Queries.GetMemberDivisionByID(c, id)
}

func (r *Repository) GetMemberDivisionByPair(memberID, divisionID uuid.NullUUID, c fiber.Ctx) (generated.MemberDivision, error) {
	return r.Queries.GetMemberDivisionByPair(c, generated.GetMemberDivisionByPairParams{
		MemberID:   memberID,
		DivisionID: divisionID,
	})
}

func (r *Repository) ListDivisionsByMemberID(memberID uuid.NullUUID, c fiber.Ctx) ([]generated.GetMemberDivisionsByMemberIDRow, error) {
	return r.Queries.GetMemberDivisionsByMemberID(c, memberID)
}

func (r *Repository) ListMembersByDivisionID(divisionID uuid.NullUUID, c fiber.Ctx) ([]generated.GetMemberDivisionsByDivisionIDRow, error) {
	return r.Queries.GetMemberDivisionsByDivisionID(c, divisionID)
}

func (r *Repository) DeleteMemberDivision(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteMemberDivision(c, id)
}
