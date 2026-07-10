package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateRoleHandler(c fiber.Ctx) error {
	var req requests.CreateRoleRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.CreateRoleService(req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Role created successfully", nil)
}

func (h *Handler) GetRoleByIDHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	role, err := h.Service.GetRoleByIDService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Role retrieved successfully", role)
}

func (h *Handler) GetAllRolesHandler(c fiber.Ctx) error {
	roles, err := h.Service.GetAllRolesService(c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Roles retrieved successfully", roles)
}

func (h *Handler) UpdateRoleHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.UpdateRoleRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.UpdateRoleService(id, req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Role updated successfully", nil)
}

func (h *Handler) DeleteRoleHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.DeleteRoleService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Role deleted successfully", nil)
}
