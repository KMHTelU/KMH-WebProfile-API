package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateBlogPostHandler(c fiber.Ctx) error {
	var req requests.CreateBlogPostRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.CreateBlogPostService(req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Blog post created successfully", nil)
}

func (h *Handler) GetBlogPostByIDHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	post, err := h.Service.GetBlogPostByIDService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Blog post retrieved successfully", post)
}

func (h *Handler) GetPaginatedBlogPostsHandler(c fiber.Ctx) error {
	limit, offset := utils.GetPaginationParams(c)
	posts, err := h.Service.GetPaginatedBlogPostsService(limit, offset, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Blog posts retrieved successfully", posts)
}

func (h *Handler) UpdateBlogPostHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.UpdateBlogPostRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.UpdateBlogPostService(id, req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Blog post updated successfully", nil)
}

func (h *Handler) DeleteBlogPostHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.DeleteBlogPostService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Blog post deleted successfully", nil)
}
