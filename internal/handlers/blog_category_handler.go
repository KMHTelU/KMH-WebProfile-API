package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateBlogCategoryHandler(c fiber.Ctx) error {
	var req requests.CreateBlogCategoryRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.CreateBlogCategoryService(req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Blog category created successfully", nil)
}

func (h *Handler) GetBlogCategoryByIDHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	category, err := h.Service.GetBlogCategoryByIDService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Blog category retrieved successfully", category)
}

func (h *Handler) GetPaginatedBlogCategoriesHandler(c fiber.Ctx) error {
	limit, offset := utils.GetPaginationParams(c)
	categories, err := h.Service.GetPaginatedBlogCategoriesService(limit, offset, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Blog categories retrieved successfully", categories)
}

func (h *Handler) UpdateBlogCategoryHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.UpdateBlogCategoryRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.UpdateBlogCategoryService(id, req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Blog category updated successfully", nil)
}

func (h *Handler) DeleteBlogCategoryHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.DeleteBlogCategoryService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Blog category deleted successfully", nil)
}
