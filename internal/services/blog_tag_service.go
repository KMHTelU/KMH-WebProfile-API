package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateBlogTagService(req requests.CreateBlogTagRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.InsertBlogTagParams{
		ID:   uuid.New(),
		Name: sql.NullString{String: req.Name, Valid: req.Name != ""},
		Slug: sql.NullString{String: req.Slug, Valid: req.Slug != ""},
	}

	if err := s.Repository.CreateBlogTag(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create blog tag")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Create Blog Tag", Valid: true},
		Entity:    sql.NullString{String: "Blog Tag", Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) GetBlogTagByIDService(id uuid.UUID, c fiber.Ctx) (generated.BlogTag, *fiber.Error) {
	tag, err := s.Repository.GetBlogTagByID(id, c)
	if err != nil {
		return generated.BlogTag{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get blog tag")
	}
	return tag, nil
}

func (s *Service) GetPaginatedBlogTagsService(limit, offset int32, c fiber.Ctx) ([]generated.BlogTag, *fiber.Error) {
	tags, err := s.Repository.ListBlogTags(generated.ListBlogTagsParams{
		Limit:  limit,
		Offset: offset,
	}, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get blog tags")
	}
	return tags, nil
}

func (s *Service) UpdateBlogTagService(id uuid.UUID, req requests.UpdateBlogTagRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.UpdateBlogTagParams{
		ID:   id,
		Name: sql.NullString{String: req.Name, Valid: req.Name != ""},
		Slug: sql.NullString{String: req.Slug, Valid: req.Slug != ""},
	}

	if err := s.Repository.UpdateBlogTag(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update blog tag")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Blog Tag", Valid: true},
		Entity:    sql.NullString{String: "Blog Tag", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) DeleteBlogTagService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteBlogTag(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete blog tag")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Blog Tag", Valid: true},
		Entity:    sql.NullString{String: "Blog Tag", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}
