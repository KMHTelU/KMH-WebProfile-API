package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateContactMessageHandler(c fiber.Ctx) error {
	var req requests.CreateContactMessageRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.CreateContactMessageService(req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Contact message submitted successfully", nil)
}

func (h *Handler) GetContactMessageByIDHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	message, err := h.Service.GetContactMessageByIDService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Contact message retrieved successfully", message)
}

func (h *Handler) GetPaginatedContactMessagesHandler(c fiber.Ctx) error {
	limit, offset := utils.GetPaginationParams(c)
	messages, err := h.Service.GetPaginatedContactMessagesService(limit, offset, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Contact messages retrieved successfully", messages)
}

func (h *Handler) MarkContactMessageAsReadHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.MarkContactMessageAsReadService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Contact message marked as read", nil)
}

func (h *Handler) DeleteContactMessageHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.DeleteContactMessageService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Contact message deleted successfully", nil)
}
