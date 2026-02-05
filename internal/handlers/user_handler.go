package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) CreateUser(c fiber.Ctx) error {
	// Implementation for creating a user
	var request requests.CreateUserRequest
	if err := c.Bind().JSON(request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}

	if err := h.Service.CreateUserService(request, c); err != nil {
		//if failed, then must check repository.InsertUser and repository.InsertLog
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	return utils.RespondWithCreated(c, "User created successfully", nil)
}
