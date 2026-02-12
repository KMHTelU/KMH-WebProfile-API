package repositories

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) InsertUser(params generated.CreateUserParams, c fiber.Ctx) error {
	_, err := r.Queries.CreateUser(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateUser(params generated.UpdateUserParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateUser(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteUser(id uuid.UUID, c fiber.Ctx) error {
	err := r.Queries.DeleteUser(c, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserByID(id uuid.UUID, c fiber.Ctx) (generated.GetUserByIDRow, error) {
	user, err := r.Queries.GetUserByID(c, id)
	if err != nil {
		return generated.GetUserByIDRow{}, err
	}
	return user, nil
}

func (r *Repository) GetAllUsers(params generated.GetUsersParams, c fiber.Ctx) ([]generated.GetUsersRow, error) {
	users, err := r.Queries.GetUsers(c, params)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByEmail(email string, c fiber.Ctx) (generated.GetUserByEmailRow, error) {
	user, err := r.Queries.GetUserByEmail(c, sql.NullString{String: email})
	if err != nil {
		return generated.GetUserByEmailRow{}, err
	}
	return user, nil
}
