package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateRoleService(req requests.CreateRoleRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.InsertRoleParams{
		ID:          uuid.New(),
		Name:        sql.NullString{String: req.Name, Valid: req.Name != ""},
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	}

	if _, err := s.Repository.CreateRole(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create role")
	}

	return nil
}

func (s *Service) GetRoleByIDService(id uuid.UUID, c fiber.Ctx) (generated.Role, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return generated.Role{}, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	role, err1 := s.Repository.GetRoleByID(id, c)
	if err1 != nil {
		return generated.Role{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get role")
	}
	return role, nil
}

func (s *Service) GetAllRolesService(c fiber.Ctx) ([]generated.Role, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	roles, err1 := s.Repository.GetAllRoles(c)
	if err1 != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get roles")
	}
	return roles, nil
}

func (s *Service) UpdateRoleService(id uuid.UUID, req requests.UpdateRoleRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.UpdateRoleParams{
		ID:          id,
		Name:        sql.NullString{String: req.Name, Valid: req.Name != ""},
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	}

	if err := s.Repository.UpdateRole(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update role")
	}
	return nil
}

func (s *Service) DeleteRoleService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteRole(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete role")
	}
	return nil
}
