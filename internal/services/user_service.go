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

	newID := uuid.New()
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	params := generated.CreateUserParams{
		ID:           newID,
		Name:         utils.NullString(req.Name),
		Email:        utils.NullString(req.Email),
		PasswordHash: utils.NullString(hashedPassword),
		RoleID:       req.RoleID,
		// User baru dianggap aktif; tanpa ini kolom is_active tetap NULL.
		IsActive: utils.NullBool(true),
	}

	if err := s.Repository.InsertUser(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Create User"),
		Entity:    utils.NullString("User with RoleID: " + req.RoleID.String()),
		EntityID:  utils.NullUUID(newID),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
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
	users, err := s.Repository.GetAllUsers(generated.GetUsersParams{
		Limit:  limit,
		Offset: offset,
	}, c)
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

	existing, err := s.Repository.GetUserByID(id, c)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	params := generated.UpdateUserParams{
		ID:           id,
		Name:         existing.Name,
		Email:        existing.Email,
		PasswordHash: existing.PasswordHash,
		RoleID:       existing.RoleID,
		IsActive:     existing.IsActive,
		// last_login_at tidak boleh direset saat update profil.
		LastLoginAt: existing.LastLoginAt,
	}

	if req.Name != "" {
		params.Name = utils.NullString(req.Name)
	}
	if req.Email != "" {
		params.Email = utils.NullString(req.Email)
	}
	if req.RoleID != uuid.Nil {
		params.RoleID = req.RoleID
	}
	if req.IsActive != nil {
		params.IsActive = sql.NullBool{Bool: *req.IsActive, Valid: true}
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
		}
		params.PasswordHash = utils.NullString(hashedPassword)
	}

	if err := s.Repository.UpdateUser(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Update User"),
		Entity:    utils.NullString("User with RoleID: " + params.RoleID.String()),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
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
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Delete User"),
		Entity:    utils.NullString("User"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
