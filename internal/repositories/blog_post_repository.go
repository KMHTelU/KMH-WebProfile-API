package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateBlogPost(params generated.InsertBlogPostParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertBlogPost(c, params)
	return err
}

func (r *Repository) GetBlogPostByID(id uuid.UUID, c fiber.Ctx) (generated.SelectBlogPostByIDRow, error) {
	return r.Queries.SelectBlogPostByID(c, id)
}

func (r *Repository) ListBlogPosts(params generated.ListBlogPostsParams, c fiber.Ctx) ([]generated.ListBlogPostsRow, error) {
	return r.Queries.ListBlogPosts(c, params)
}

func (r *Repository) ListBlogPostsByCategory(params generated.ListBlogPostsByCategoryParams, c fiber.Ctx) ([]generated.ListBlogPostsByCategoryRow, error) {
	return r.Queries.ListBlogPostsByCategory(c, params)
}

func (r *Repository) UpdateBlogPost(params generated.UpdateBlogPostParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateBlogPost(c, params)
	return err
}

func (r *Repository) DeleteBlogPost(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteBlogPost(c, id)
}

func (r *Repository) InsertBlogPostTag(params generated.InsertBlogPostTagParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertBlogPostTag(c, params)
	return err
}

func (r *Repository) DeleteBlogPostTag(params generated.DeleteBlogPostTagParams, c fiber.Ctx) error {
	return r.Queries.DeleteBlogPostTag(c, params)
}

func (r *Repository) ListTagsByBlogPostID(params generated.ListTagsByBlogPostIDParams, c fiber.Ctx) ([]generated.BlogTag, error) {
	return r.Queries.ListTagsByBlogPostID(c, params)
}
