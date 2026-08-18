package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreateGallery(params generated.InsertGalleryParams, c fiber.Ctx) (generated.Gallery, error) {
	return r.Queries.InsertGallery(c, params)
}

func (r *Repository) GetGalleryByID(id uuid.UUID, c fiber.Ctx) (generated.SelectGalleryByIDRow, error) {
	return r.Queries.SelectGalleryByID(c, id)
}

func (r *Repository) ListGalleries(params generated.SelectAllGalleriesParams, c fiber.Ctx) ([]generated.SelectAllGalleriesRow, error) {
	return r.Queries.SelectAllGalleries(c, params)
}

func (r *Repository) ListGalleryCategories(c fiber.Ctx) ([]generated.SelectGalleryCategoriesRow, error) {
	return r.Queries.SelectGalleryCategories(c)
}

func (r *Repository) GetGalleryItemByID(id uuid.UUID, c fiber.Ctx) (generated.GalleryItem, error) {
	return r.Queries.GetGalleryItemByID(c, id)
}

func (r *Repository) UpdateGalleryItem(params generated.UpdateGalleryItemParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateGalleryItem(c, params)
	return err
}

func (r *Repository) UpdateGallery(params generated.UpdateGalleryParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateGallery(c, params)
	return err
}

func (r *Repository) DeleteGallery(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteGallery(c, id)
}

func (r *Repository) InsertGalleryItem(params generated.InsertGalleryItemParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertGalleryItem(c, params)
	return err
}

func (r *Repository) DeleteGalleryItem(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteGalleryItem(c, id)
}

func (r *Repository) DeleteGalleryItemsByGalleryID(galleryID uuid.NullUUID, c fiber.Ctx) error {
	return r.Queries.DeleteGalleryItemsByGalleryID(c, galleryID)
}

func (r *Repository) ListGalleryItemsByGalleryID(galleryID uuid.NullUUID, c fiber.Ctx) ([]generated.SelectGalleryItemsByGalleryIDRow, error) {
	return r.Queries.SelectGalleryItemsByGalleryID(c, galleryID)
}
