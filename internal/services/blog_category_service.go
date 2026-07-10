package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateBlogCategoryService(req requests.CreateBlogCategoryRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.InsertBlogCategoryParams{
		ID:   uuid.New(),
		Name: sql.NullString{String: req.Name, Valid: req.Name != ""},
		Slug: sql.NullString{String: req.Slug, Valid: req.Slug != ""},
	}

	if err := s.Repository.CreateBlogCategory(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create blog category")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Create Blog Category", Valid: true},
		Entity:    sql.NullString{String: "Blog Category", Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) GetBlogCategoryByIDService(id uuid.UUID, c fiber.Ctx) (generated.BlogCategory, *fiber.Error) {
	category, err := s.Repository.GetBlogCategoryByID(id, c)
	if err != nil {
		return generated.BlogCategory{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get blog category")
	}
	return category, nil
}

func (s *Service) GetPaginatedBlogCategoriesService(limit, offset int32, c fiber.Ctx) ([]generated.BlogCategory, *fiber.Error) {
	categories, err := s.Repository.ListBlogCategories(generated.ListBlogCategoriesParams{
		Limit:  limit,
		Offset: offset,
	}, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get blog categories")
	}
	return categories, nil
}

func (s *Service) UpdateBlogCategoryService(id uuid.UUID, req requests.UpdateBlogCategoryRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.UpdateBlogCategoryParams{
		ID:   id,
		Name: sql.NullString{String: req.Name, Valid: req.Name != ""},
		Slug: sql.NullString{String: req.Slug, Valid: req.Slug != ""},
	}

	if err := s.Repository.UpdateBlogCategory(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update blog category")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Blog Category", Valid: true},
		Entity:    sql.NullString{String: "Blog Category", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) DeleteBlogCategoryService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteBlogCategory(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete blog category")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Blog Category", Valid: true},
		Entity:    sql.NullString{String: "Blog Category", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}
