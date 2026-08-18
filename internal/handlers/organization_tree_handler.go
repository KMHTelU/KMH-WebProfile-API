package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) GetOrganizationTreeHandler(c fiber.Ctx) error {
	tree, err := h.Service.GetOrganizationTreeService(c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Organization tree retrieved successfully", tree)
}
