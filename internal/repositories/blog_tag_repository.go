package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateBlogTag(params generated.InsertBlogTagParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertBlogTag(c, params)
	return err
}

func (r *Repository) GetBlogTagByID(id uuid.UUID, c fiber.Ctx) (generated.BlogTag, error) {
	return r.Queries.SelectBlogTagByID(c, id)
}

func (r *Repository) ListBlogTags(params generated.ListBlogTagsParams, c fiber.Ctx) ([]generated.BlogTag, error) {
	return r.Queries.ListBlogTags(c, params)
}

func (r *Repository) UpdateBlogTag(params generated.UpdateBlogTagParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateBlogTag(c, params)
	return err
}

func (r *Repository) DeleteBlogTag(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteBlogTag(c, id)
}
