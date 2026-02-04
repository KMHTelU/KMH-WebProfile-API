package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateUserService(req requests.CreateUserRequest, c fiber.Ctx) error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return utils.RespondWithError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	newId := uuid.New()
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to hash password")
	}
	var params generated.CreateUserParams = generated.CreateUserParams{
		ID:           newId,
		Name:         sql.NullString{String: req.Name},
		Email:        sql.NullString{String: req.Email},
		PasswordHash: sql.NullString{String: hashedPassword},
		RoleID:       req.RoleID,
	}

	if err := s.Repository.InsertUser(params, c); err != nil {
		return err
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
		return err
	}
	return nil
}
