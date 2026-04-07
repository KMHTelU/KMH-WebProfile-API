package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateMember(params generated.InsertMemberParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertMember(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMemberByID(id uuid.UUID, c fiber.Ctx) (generated.GetMemberByIDRow, error) {
	member, err := r.Queries.GetMemberByID(c, id)
	if err != nil {
		return generated.GetMemberByIDRow{}, err
	}
	return member, nil
}

func (r *Repository) GetAllMembers(params generated.GetAllMembersParams, c fiber.Ctx) ([]generated.GetAllMembersRow, error) {
	members, err := r.Queries.GetAllMembers(c, params)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *Repository) UpdateMember(params generated.UpdateMemberParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateMember(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteMember(id uuid.UUID, c fiber.Ctx) error {
	err := r.Queries.DeleteMember(c, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateMemberPhoto(params generated.UpdateMemberPhotoParams, c fiber.Ctx) error {
	err := r.Queries.UpdateMemberPhoto(c, params)
	if err != nil {
		return err
	}
	return nil
}
