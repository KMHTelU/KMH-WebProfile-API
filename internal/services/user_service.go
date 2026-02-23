package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateUserService(req requests.CreateUserRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	newId := uuid.New()
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}
	var params generated.CreateUserParams = generated.CreateUserParams{
		ID:           newId,
		Name:         sql.NullString{String: req.Name},
		Email:        sql.NullString{String: req.Email},
		PasswordHash: sql.NullString{String: hashedPassword},
		RoleID:       req.RoleID,
	}

	if err := s.Repository.InsertUser(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Create User"},
		Entity:    sql.NullString{String: "User with RoleID: " + req.RoleID.String()},
		EntityID:  uuid.NullUUID{UUID: newId},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) GetAllUsersService(limit, offset int32, c fiber.Ctx) ([]generated.GetUsersRow, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	userParam := generated.GetUsersParams{
		Limit:  limit,
		Offset: offset,
	}
	users, err := s.Repository.GetAllUsers(userParam, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get users")
	}
	return users, nil
}

func (s *Service) GetUserByIDService(id uuid.UUID, c fiber.Ctx) (generated.GetUserByIDRow, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return generated.GetUserByIDRow{}, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	user, err := s.Repository.GetUserByID(id, c)
	if err != nil {
		return generated.GetUserByIDRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get user")
	}
	return user, nil
}

func (s *Service) UpdateUserService(id uuid.UUID, req requests.UpdateUserRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}
	var params generated.UpdateUserParams = generated.UpdateUserParams{
		ID:           id,
		Name:         sql.NullString{String: req.Name},
		Email:        sql.NullString{String: req.Email},
		PasswordHash: sql.NullString{String: hashedPassword},
		RoleID:       req.RoleID,
	}
	if err := s.Repository.UpdateUser(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Update User"},
		Entity:    sql.NullString{String: "User with RoleID: " + req.RoleID.String()},
		EntityID:  uuid.NullUUID{UUID: id},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) DeleteUserService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	if err := s.Repository.DeleteUser(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Delete User"},
		Entity:    sql.NullString{String: "User"},
		EntityID:  uuid.NullUUID{UUID: id},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
