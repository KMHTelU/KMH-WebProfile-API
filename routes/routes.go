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
	// Unprotected routes
	api := app.Group("/api")
	api.Get("/status", func(c fiber.Ctx) error {
		return utils.RespondWithOK(c, "KMH WebProfile API is running", nil)
	})
	api.Get("/homepage-banners", func(c fiber.Ctx) error {
		return r.Handler.GetHomepageBannersHandler(c)
	})

	// Authentication routes
	api.Post("/login", func(c fiber.Ctx) error {
		return r.Handler.AuthenticateUser(c)
	})
	api.Post("/refresh", func(c fiber.Ctx) error {
		return r.Handler.RefreshToken(c)
	})

	// Protected routes
	protected := api.Group("/protected")

	user := protected.Group("/user")
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

	homepageBanner := protected.Group("/homepage-banners")
	homepageBanner.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateHomepageBannerHandler(c)
	})
	homepageBanner.Get("", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedHomepageBannersHandler(c)
	})
	homepageBanner.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteHomepageBannerHandler(c)
	})
}
