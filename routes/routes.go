package routes

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/handlers"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/status", func(c fiber.Ctx) error {
		return utils.RespondWithOK(c, "KMH WebProfile API is running", nil)
	})
	// Additional routes can be added here
	user := api.Group("/user")
	user.Post("/", func(c fiber.Ctx) error {
		return handlers.CreateUser(c)
	})
}
