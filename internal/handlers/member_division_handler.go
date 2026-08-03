package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateMemberDivisionHandler(c fiber.Ctx) error {
	var req requests.CreateMemberDivisionRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	id, err := h.Service.CreateMemberDivisionService(req, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithCreated(c, "Member assigned to division successfully", fiber.Map{"id": id})
}

func (h *Handler) UpdateMemberDivisionHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.UpdateMemberDivisionRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	if err := h.Service.UpdateMemberDivisionService(id, req, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Member division updated successfully", nil)
}

func (h *Handler) DeleteMemberDivisionHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if err := h.Service.DeleteMemberDivisionService(id, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Member division deleted successfully", nil)
}

func (h *Handler) GetDivisionsByMemberHandler(c fiber.Ctx) error {
	memberID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid Member ID")
	}
	rows, ferr := h.Service.GetDivisionsByMemberService(memberID, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "Divisions of member retrieved successfully", rows)
}

func (h *Handler) GetMembersByDivisionHandler(c fiber.Ctx) error {
	divisionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid Division ID")
	}
	rows, ferr := h.Service.GetMembersByDivisionService(divisionID, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "Members of division retrieved successfully", rows)
}
