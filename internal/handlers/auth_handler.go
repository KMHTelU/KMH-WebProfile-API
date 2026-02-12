package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) AuthenticateUser(c fiber.Ctx) error {
	// Implementation for authenticating a user
	var request requests.AuthenticateUserRequest
	if err := c.Bind().JSON(&request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	auth, err := h.Service.AuthenticateUserService(request, c)
	if err != nil {
		// if failed, the error is already handled in service
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "User authenticated successfully", auth)
}

func (h *Handler) RefreshToken(c fiber.Ctx) error {
	// Implementation for refreshing JWT token
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind().JSON(&request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	auth, err := h.Service.RefreshTokenService(request.RefreshToken, c)
	if err != nil {
		// if failed, the error is already handled in service
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Token refreshed successfully", auth)
}
