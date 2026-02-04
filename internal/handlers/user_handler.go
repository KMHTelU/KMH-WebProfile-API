package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

func CreateUser(c fiber.Ctx) error {
	// Implementation for creating a user
	var request requests.CreateUserRequest
	if err := c.Bind().JSON(request); err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}

	return utils.RespondWithCreated(c, "User created successfully", nil)
}
