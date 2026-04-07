package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateMemberHandler(c fiber.Ctx) error {
	var req requests.CreateMemberRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.CreateMemberService(req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Member created successfully", nil)
}

func (h *Handler) GetMemberByIDHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	member, err := h.Service.GetMemberByIDService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Member retrieved successfully", member)
}

func (h *Handler) GetPaginatedAllMembersHandler(c fiber.Ctx) error {
	limit, offset := utils.GetPaginationParams(c)
	members, err := h.Service.GetPaginatedAllMembersService(limit, offset, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Members retrieved successfully", members)
}

func (h *Handler) UpdateMemberHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.UpdateMemberRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.UpdateMemberService(id, req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Member updated successfully", nil)
}

func (h *Handler) DeleteMemberHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.DeleteMemberService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Member deleted successfully", nil)
}

func (h *Handler) UploadAndUpdateMemberPhotoHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.UploadMemberPhotoService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Member photo uploaded and updated successfully", nil)
}
