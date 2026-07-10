package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateBlogCategory(params generated.InsertBlogCategoryParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertBlogCategory(c, params)
	return err
}

func (r *Repository) GetBlogCategoryByID(id uuid.UUID, c fiber.Ctx) (generated.BlogCategory, error) {
	return r.Queries.GetBlogCategoryByID(c, id)
}

func (r *Repository) ListBlogCategories(params generated.ListBlogCategoriesParams, c fiber.Ctx) ([]generated.BlogCategory, error) {
	return r.Queries.ListBlogCategories(c, params)
}

func (r *Repository) UpdateBlogCategory(params generated.UpdateBlogCategoryParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateBlogCategory(c, params)
	return err
}

func (r *Repository) DeleteBlogCategory(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteBlogCategory(c, id)
}
