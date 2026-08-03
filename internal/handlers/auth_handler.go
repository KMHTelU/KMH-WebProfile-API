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

func (h *Handler) ForgotPasswordHandler(c fiber.Ctx) error {
	var request requests.ForgotPasswordRequest
	if err := c.Bind().JSON(&request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	if err := h.Service.ForgotPasswordService(request, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	// Pesan sengaja netral supaya tidak membocorkan email mana yang terdaftar.
	return utils.RespondWithOK(c, "If the email is registered, a password reset link has been sent", nil)
}

func (h *Handler) ResetPasswordHandler(c fiber.Ctx) error {
	var request requests.ResetPasswordRequest
	if err := c.Bind().JSON(&request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	if err := h.Service.ResetPasswordService(request, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Password reset successfully", nil)
}

func (h *Handler) ChangePasswordHandler(c fiber.Ctx) error {
	var request requests.ChangePasswordRequest
	if err := c.Bind().JSON(&request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	if err := h.Service.ChangePasswordService(request, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Password changed successfully", nil)
}
