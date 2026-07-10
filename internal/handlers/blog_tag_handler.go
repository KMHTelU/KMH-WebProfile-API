package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateBlogTagHandler(c fiber.Ctx) error {
	var req requests.CreateBlogTagRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.CreateBlogTagService(req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Blog tag created successfully", nil)
}

func (h *Handler) GetBlogTagByIDHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	tag, err := h.Service.GetBlogTagByIDService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Blog tag retrieved successfully", tag)
}

func (h *Handler) GetPaginatedBlogTagsHandler(c fiber.Ctx) error {
	limit, offset := utils.GetPaginationParams(c)
	tags, err := h.Service.GetPaginatedBlogTagsService(limit, offset, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Blog tags retrieved successfully", tags)
}

func (h *Handler) UpdateBlogTagHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.UpdateBlogTagRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.UpdateBlogTagService(id, req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Blog tag updated successfully", nil)
}

func (h *Handler) DeleteBlogTagHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.DeleteBlogTagService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Blog tag deleted successfully", nil)
}
