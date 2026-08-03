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
	api.Get("/organization-profile/:id", func(c fiber.Ctx) error {
		return r.Handler.GetOrganizationProfileHandler(c)
	})
	api.Get("/blog-categories", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedBlogCategoriesHandler(c)
	})
	api.Get("/blog-categories/:id", func(c fiber.Ctx) error {
		return r.Handler.GetBlogCategoryByIDHandler(c)
	})
	api.Get("/blog-tags", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedBlogTagsHandler(c)
	})
	api.Get("/blog-tags/:id", func(c fiber.Ctx) error {
		return r.Handler.GetBlogTagByIDHandler(c)
	})
	api.Get("/blog-posts", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedBlogPostsHandler(c)
	})
	api.Get("/blog-posts/:id", func(c fiber.Ctx) error {
		return r.Handler.GetBlogPostByIDHandler(c)
	})
	api.Get("/events", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedEventsHandler(c)
	})
	api.Get("/events/:id", func(c fiber.Ctx) error {
		return r.Handler.GetEventByIDHandler(c)
	})
	api.Get("/galleries", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedGalleriesHandler(c)
	})
	// Didaftarkan sebelum /galleries/:id agar "categories" tidak tertangkap
	// sebagai nilai parameter id.
	api.Get("/galleries/categories", func(c fiber.Ctx) error {
		return r.Handler.GetGalleryCategoriesHandler(c)
	})
	api.Get("/galleries/:id", func(c fiber.Ctx) error {
		return r.Handler.GetGalleryByIDHandler(c)
	})
	api.Post("/contact-messages", func(c fiber.Ctx) error {
		return r.Handler.CreateContactMessageHandler(c)
	})

	// Divisions (public read-only)
	api.Get("/divisions", func(c fiber.Ctx) error {
		return r.Handler.GetAllDivisionsHandler(c)
	})
	api.Get("/divisions/:id", func(c fiber.Ctx) error {
		return r.Handler.GetDivisionByIDHandler(c)
	})
	api.Get("/divisions/:id/members", func(c fiber.Ctx) error {
		return r.Handler.GetMembersByDivisionHandler(c)
	})

	// Members (public read-only)
	api.Get("/members", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedAllMembersHandler(c)
	})
	api.Get("/members/:id", func(c fiber.Ctx) error {
		return r.Handler.GetMemberByIDHandler(c)
	})
	api.Get("/members/:id/divisions", func(c fiber.Ctx) error {
		return r.Handler.GetDivisionsByMemberHandler(c)
	})

	// Authentication routes
	api.Post("/login", func(c fiber.Ctx) error {
		return r.Handler.AuthenticateUser(c)
	})
	api.Post("/refresh", func(c fiber.Ctx) error {
		return r.Handler.RefreshToken(c)
	})
	api.Post("/forgot-password", func(c fiber.Ctx) error {
		return r.Handler.ForgotPasswordHandler(c)
	})
	api.Post("/reset-password", func(c fiber.Ctx) error {
		return r.Handler.ResetPasswordHandler(c)
	})

	// Protected routes
	protected := api.Group("/protected")

	protected.Post("/auth/change-password", func(c fiber.Ctx) error {
		return r.Handler.ChangePasswordHandler(c)
	})

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
	// Rute /bulk didaftarkan lebih dahulu agar tidak tertangkap oleh /:id.
	homepageBanner.Post("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkCreateHomepageBannersHandler(c)
	})
	homepageBanner.Put("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkUpdateHomepageBannersHandler(c)
	})
	homepageBanner.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteHomepageBannerHandler(c)
	})

	member := protected.Group("/members")
	member.Post("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkCreateMembersHandler(c)
	})
	member.Put("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkUpdateMembersHandler(c)
	})
	member.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateMemberHandler(c)
	})
	member.Get("/:id", func(c fiber.Ctx) error {
		return r.Handler.GetMemberByIDHandler(c)
	})
	member.Get("", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedAllMembersHandler(c)
	})
	member.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateMemberHandler(c)
	})
	member.Post("/:id/upload", func(c fiber.Ctx) error {
		return r.Handler.UploadAndUpdateMemberPhotoHandler(c)
	})
	member.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteMemberHandler(c)
	})

	division := protected.Group("/divisions")
	division.Post("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkCreateDivisionsHandler(c)
	})
	division.Put("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkUpdateDivisionsHandler(c)
	})
	division.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateDivisionHandler(c)
	})
	division.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateDivisionHandler(c)
	})
	division.Post("/:id/upload", func(c fiber.Ctx) error {
		return r.Handler.UploadAndUpdateDivisionIconHandler(c)
	})
	division.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteDivisionHandler(c)
	})

	orgProfile := protected.Group("/organization-profile")
	orgProfile.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateOrganizationProfileHandler(c)
	})
	orgProfile.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateOrganizationProfileHandler(c)
	})
	orgProfile.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteOrganizationProfileHandler(c)
	})
	orgProfile.Post("/:id/upload", func(c fiber.Ctx) error {
		return r.Handler.UploadAndUpdateOrganizationProfileLogoHandler(c)
	})

	blogCategory := protected.Group("/blog-categories")
	blogCategory.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateBlogCategoryHandler(c)
	})
	blogCategory.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateBlogCategoryHandler(c)
	})
	blogCategory.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteBlogCategoryHandler(c)
	})

	blogTag := protected.Group("/blog-tags")
	blogTag.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateBlogTagHandler(c)
	})
	blogTag.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateBlogTagHandler(c)
	})
	blogTag.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteBlogTagHandler(c)
	})

	blogPost := protected.Group("/blog-posts")
	blogPost.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateBlogPostHandler(c)
	})
	blogPost.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateBlogPostHandler(c)
	})
	blogPost.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteBlogPostHandler(c)
	})

	event := protected.Group("/events")
	event.Post("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkCreateEventsHandler(c)
	})
	event.Put("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkUpdateEventsHandler(c)
	})
	event.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateEventHandler(c)
	})
	event.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateEventHandler(c)
	})
	event.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteEventHandler(c)
	})

	gallery := protected.Group("/galleries")
	gallery.Post("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkCreateGalleriesHandler(c)
	})
	gallery.Put("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkUpdateGalleriesHandler(c)
	})
	gallery.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateGalleryHandler(c)
	})
	gallery.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateGalleryHandler(c)
	})
	gallery.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteGalleryHandler(c)
	})
	gallery.Post("/:id/items", func(c fiber.Ctx) error {
		return r.Handler.AddGalleryItemHandler(c)
	})
	gallery.Delete("/:id/items/:itemId", func(c fiber.Ctx) error {
		return r.Handler.DeleteGalleryItemHandler(c)
	})

	// Gallery item punya grup tersendiri karena operasi bulk menyertakan
	// gallery_id di dalam body, bukan di path.
	galleryItem := protected.Group("/gallery-items")
	galleryItem.Post("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkCreateGalleryItemsHandler(c)
	})
	galleryItem.Put("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkUpdateGalleryItemsHandler(c)
	})

	memberDivision := protected.Group("/member-divisions")
	memberDivision.Post("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkCreateMemberDivisionsHandler(c)
	})
	memberDivision.Put("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkUpdateMemberDivisionsHandler(c)
	})
	memberDivision.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateMemberDivisionHandler(c)
	})
	memberDivision.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateMemberDivisionHandler(c)
	})
	memberDivision.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteMemberDivisionHandler(c)
	})

	contactMessage := protected.Group("/contact-messages")
	contactMessage.Get("", func(c fiber.Ctx) error {
		return r.Handler.GetPaginatedContactMessagesHandler(c)
	})
	contactMessage.Get("/:id", func(c fiber.Ctx) error {
		return r.Handler.GetContactMessageByIDHandler(c)
	})
	contactMessage.Patch("/:id/read", func(c fiber.Ctx) error {
		return r.Handler.MarkContactMessageAsReadHandler(c)
	})
	contactMessage.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteContactMessageHandler(c)
	})

	// Media (upload generik → mengembalikan media_id untuk blog/event/gallery)
	media := protected.Group("/media")
	media.Post("", func(c fiber.Ctx) error {
		return r.Handler.UploadMediaHandler(c)
	})

	role := protected.Group("/roles")
	role.Get("", func(c fiber.Ctx) error {
		return r.Handler.GetAllRolesHandler(c)
	})
	role.Post("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkCreateRolesHandler(c)
	})
	role.Put("/bulk", func(c fiber.Ctx) error {
		return r.Handler.BulkUpdateRolesHandler(c)
	})
	role.Get("/:id", func(c fiber.Ctx) error {
		return r.Handler.GetRoleByIDHandler(c)
	})
	role.Post("", func(c fiber.Ctx) error {
		return r.Handler.CreateRoleHandler(c)
	})
	role.Put("/:id", func(c fiber.Ctx) error {
		return r.Handler.UpdateRoleHandler(c)
	})
	role.Delete("/:id", func(c fiber.Ctx) error {
		return r.Handler.DeleteRoleHandler(c)
	})

	// Import berkas Excel/CSV. Nilai :entity yang didukung adalah members,
	// divisions, member-divisions, events, dan roles.
	importGroup := protected.Group("/import")
	importGroup.Get("/:entity/template", func(c fiber.Ctx) error {
		return r.Handler.DownloadImportTemplateHandler(c)
	})
	importGroup.Post("/:entity/preview", func(c fiber.Ctx) error {
		return r.Handler.PreviewImportHandler(c)
	})
	importGroup.Post("/:entity", func(c fiber.Ctx) error {
		return r.Handler.ImportHandler(c)
	})
}
