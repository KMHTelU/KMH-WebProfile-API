package routes

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/handlers"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

type Routes struct {
	Handler *handlers.Handler
}

func InitializeRoutes(handler *handlers.Handler) *Routes {
	return &Routes{
		Handler: handler,
	}
}

func (r *Routes) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/status", func(c fiber.Ctx) error {
		return utils.RespondWithOK(c, "KMH WebProfile API is running", nil)
	})
	// Additional routes can be added here
	user := api.Group("/user")
	user.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateUser(c)
	})
	user.Get("", func(c fiber.Ctx) error {
		return r.Handler.GetAllUsers(c)
	})
	user.Get("/:id", func(c fiber.Ctx) error {
		return r.Handler.GetUserByID(c)
	})
	user.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateUser(c)
	})
	user.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteUser(c)
	})

	auth := api.Group("/auth")
	auth.Post("/login", func(c fiber.Ctx) error {
		return r.Handler.AuthenticateUser(c)
	})
	auth.Post("/refresh", func(c fiber.Ctx) error {
		return r.Handler.RefreshToken(c)
	})
}
