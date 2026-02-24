package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateHomepageBannerHandler(c fiber.Ctx) error {
	var req requests.HomepageBannerRequest
	if err := c.Bind().Form(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}

	if err := h.Service.CreateHomepageBannerService(req, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithCreated(c, "Homepage banner created successfully", nil)
}

func (h *Handler) GetHomepageBannersHandler(c fiber.Ctx) error {
	banners, err := h.Service.GetHomepageBannersService(c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Homepage banners retrieved successfully", banners)
}
func (h *Handler) GetPaginatedHomepageBannersHandler(c fiber.Ctx) error {
	limit, offset := utils.GetPaginationParams(c)
	banners, err := h.Service.GetPaginatedHomepageBannersService(limit, offset, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Homepage banners retrieved successfully", banners)
}

func (h *Handler) DeleteHomepageBannerHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if err := h.Service.DeleteHomepageBannerService(id, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Homepage banner deleted successfully", nil)
}
