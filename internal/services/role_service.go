package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateRoleService(req requests.CreateRoleRequest, c fiber.Ctx) *fiber.Error {
	if _, ferr := s.requireAuth(c); ferr != nil {
		return ferr
	}
	if _, ferr := s.createRole(req, c); ferr != nil {
		return ferr
	}
	return nil
}

func (s *Service) GetRoleByIDService(id uuid.UUID, c fiber.Ctx) (generated.Role, *fiber.Error) {
	if _, ferr := s.requireAuth(c); ferr != nil {
		return generated.Role{}, ferr
	}
	role, err := s.Repository.GetRoleByID(id, c)
	if err != nil {
		return generated.Role{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get role")
	}
	return role, nil
}

func (s *Service) GetAllRolesService(c fiber.Ctx) ([]generated.Role, *fiber.Error) {
	if _, ferr := s.requireAuth(c); ferr != nil {
		return nil, ferr
	}
	roles, err := s.Repository.GetAllRoles(c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get roles")
	}
	return roles, nil
}

func (s *Service) UpdateRoleService(id uuid.UUID, req requests.UpdateRoleRequest, c fiber.Ctx) *fiber.Error {
	if _, ferr := s.requireAuth(c); ferr != nil {
		return ferr
	}
	return s.updateRole(id, req, c)
}

func (s *Service) DeleteRoleService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	if _, ferr := s.requireAuth(c); ferr != nil {
		return ferr
	}

	if err := s.Repository.DeleteRole(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete role")
	}
	return nil
}
