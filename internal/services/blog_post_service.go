package services

import (
	"database/sql"
	"time"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateBlogPostService(req requests.CreateBlogPostRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	var pubAt sql.NullTime
	if req.Status == "PUBLISHED" {
		pubAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	params := generated.InsertBlogPostParams{
		ID:              uuid.New(),
		Title:           sql.NullString{String: req.Title, Valid: req.Title != ""},
		Slug:            sql.NullString{String: req.Slug, Valid: req.Slug != ""},
		Excerpt:         sql.NullString{String: req.Excerpt, Valid: req.Excerpt != ""},
		Content:         sql.NullString{String: req.Content, Valid: req.Content != ""},
		CategoryID:      uuid.NullUUID{UUID: req.CategoryID, Valid: req.CategoryID != uuid.Nil},
		FeaturedMediaID: uuid.NullUUID{UUID: req.FeaturedMediaID, Valid: req.FeaturedMediaID != uuid.Nil},
		AuthorID:        uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Status:          sql.NullString{String: req.Status, Valid: req.Status != ""},
		PublishedAt:     pubAt,
	}

	if err := s.Repository.CreateBlogPost(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create blog post")
	}

	for _, tagID := range req.TagIDs {
		_ = s.Repository.InsertBlogPostTag(generated.InsertBlogPostTagParams{
			PostID: params.ID,
			TagID:  tagID,
		}, c)
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Create Blog Post", Valid: true},
		Entity:    sql.NullString{String: "Blog Post", Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) GetBlogPostByIDService(id uuid.UUID, c fiber.Ctx) (map[string]interface{}, *fiber.Error) {
	post, err := s.Repository.GetBlogPostByID(id, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get blog post")
	}
	tags, _ := s.Repository.ListTagsByBlogPostID(generated.ListTagsByBlogPostIDParams{
		PostID: id,
		Limit:  100,
		Offset: 0,
	}, c)

	return map[string]interface{}{
		"post": post,
		"tags": tags,
	}, nil
}

func (s *Service) GetPaginatedBlogPostsService(limit, offset int32, c fiber.Ctx) ([]generated.ListBlogPostsRow, *fiber.Error) {
	posts, err := s.Repository.ListBlogPosts(generated.ListBlogPostsParams{
		Limit:  limit,
		Offset: offset,
	}, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get blog posts")
	}
	return posts, nil
}

func (s *Service) UpdateBlogPostService(id uuid.UUID, req requests.UpdateBlogPostRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	existingPost, err := s.Repository.GetBlogPostByID(id, c)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Blog post not found")
	}

	var pubAt sql.NullTime = existingPost.PublishedAt
	if existingPost.Status.String != "PUBLISHED" && req.Status == "PUBLISHED" {
		pubAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	params := generated.UpdateBlogPostParams{
		ID:              id,
		Title:           sql.NullString{String: req.Title, Valid: req.Title != ""},
		Slug:            sql.NullString{String: req.Slug, Valid: req.Slug != ""},
		Excerpt:         sql.NullString{String: req.Excerpt, Valid: req.Excerpt != ""},
		Content:         sql.NullString{String: req.Content, Valid: req.Content != ""},
		CategoryID:      uuid.NullUUID{UUID: req.CategoryID, Valid: req.CategoryID != uuid.Nil},
		FeaturedMediaID: uuid.NullUUID{UUID: req.FeaturedMediaID, Valid: req.FeaturedMediaID != uuid.Nil},
		AuthorID:        uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Status:          sql.NullString{String: req.Status, Valid: req.Status != ""},
		PublishedAt:     pubAt,
	}

	if err := s.Repository.UpdateBlogPost(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update blog post")
	}

	// Update tags: naive approach delete all and insert new
	existingTags, _ := s.Repository.ListTagsByBlogPostID(generated.ListTagsByBlogPostIDParams{
		PostID: id, Limit: 100, Offset: 0,
	}, c)
	for _, t := range existingTags {
		_ = s.Repository.DeleteBlogPostTag(generated.DeleteBlogPostTagParams{
			PostID: id,
			TagID:  t.ID,
		}, c)
	}
	for _, tagID := range req.TagIDs {
		_ = s.Repository.InsertBlogPostTag(generated.InsertBlogPostTagParams{
			PostID: id,
			TagID:  tagID,
		}, c)
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Blog Post", Valid: true},
		Entity:    sql.NullString{String: "Blog Post", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) DeleteBlogPostService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteBlogPost(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete blog post")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Blog Post", Valid: true},
		Entity:    sql.NullString{String: "Blog Post", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}
