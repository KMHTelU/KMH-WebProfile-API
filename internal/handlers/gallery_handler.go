package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateGalleryHandler(c fiber.Ctx) error {
	var req requests.CreateGalleryRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	gallery, err := h.Service.CreateGalleryService(req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Gallery created successfully", gallery)
}

func (h *Handler) GetGalleryByIDHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	galleryData, err := h.Service.GetGalleryByIDService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Gallery retrieved successfully", galleryData)
}

func (h *Handler) GetPaginatedGalleriesHandler(c fiber.Ctx) error {
	limit, offset := utils.GetPaginationParams(c)
	galleries, err := h.Service.GetPaginatedGalleriesService(limit, offset, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Galleries retrieved successfully", galleries)
}

func (h *Handler) UpdateGalleryHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.UpdateGalleryRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.UpdateGalleryService(id, req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Gallery updated successfully", nil)
}

func (h *Handler) DeleteGalleryHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.DeleteGalleryService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Gallery deleted successfully", nil)
}

func (h *Handler) AddGalleryItemHandler(c fiber.Ctx) error {
	galleryID, errP := uuid.Parse(c.Params("id"))
	if errP != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid Gallery ID")
	}
	var req requests.AddGalleryItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.AddGalleryItemService(galleryID, req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Gallery item added successfully", nil)
}

func (h *Handler) DeleteGalleryItemHandler(c fiber.Ctx) error {
	itemID, errP := uuid.Parse(c.Params("itemId"))
	if errP != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid Item ID")
	}
	err := h.Service.DeleteGalleryItemService(itemID, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Gallery item deleted successfully", nil)
}
