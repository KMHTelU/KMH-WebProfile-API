package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateGalleryService(req requests.CreateGalleryRequest, c fiber.Ctx) (generated.Gallery, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return generated.Gallery{}, ferr
	}

	gallery, ferr := s.createGallery(req, c)
	if ferr != nil {
		return generated.Gallery{}, ferr
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Create Gallery"),
		Entity:    utils.NullString("Gallery"),
		EntityID:  utils.NullUUID(gallery.ID),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return gallery, nil
}

func (s *Service) GetGalleryByIDService(id uuid.UUID, c fiber.Ctx) (map[string]interface{}, *fiber.Error) {
	gallery, err := s.Repository.GetGalleryByID(id, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get gallery")
	}

	items, _ := s.Repository.ListGalleryItemsByGalleryID(uuid.NullUUID{UUID: id, Valid: true}, c)

	return map[string]interface{}{
		"gallery": gallery,
		"items":   items,
	}, nil
}

// GalleryFilter menampung filter opsional pada daftar gallery. Field yang tidak
// diisi berarti tidak menyaring apa pun.
type GalleryFilter struct {
	EventID  uuid.NullUUID
	IsPublic sql.NullBool
}

func (s *Service) GetPaginatedGalleriesService(limit, offset int32, filter GalleryFilter, c fiber.Ctx) ([]generated.SelectAllGalleriesRow, *fiber.Error) {
	galleries, err := s.Repository.ListGalleries(generated.SelectAllGalleriesParams{
		Limit:    limit,
		Offset:   offset,
		EventID:  filter.EventID,
		IsPublic: filter.IsPublic,
	}, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get galleries")
	}
	return galleries, nil
}

// GetGalleryCategoriesService mengembalikan event yang memiliki gallery.
// Frontend memakai daftar ini sebagai pilihan filter pada halaman galeri.
func (s *Service) GetGalleryCategoriesService(c fiber.Ctx) ([]generated.SelectGalleryCategoriesRow, *fiber.Error) {
	categories, err := s.Repository.ListGalleryCategories(c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get gallery categories")
	}
	return categories, nil
}

func (s *Service) UpdateGalleryService(id uuid.UUID, req requests.UpdateGalleryRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	if ferr := s.updateGallery(id, req, c); ferr != nil {
		return ferr
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Update Gallery"),
		Entity:    utils.NullString("Gallery"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

func (s *Service) DeleteGalleryService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	// Item gallery harus dihapus lebih dulu karena gallery_items punya foreign
	// key ke galleries tanpa ON DELETE CASCADE; kalau tidak, delete gagal.
	if err := s.Repository.DeleteGalleryItemsByGalleryID(uuid.NullUUID{UUID: id, Valid: true}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete gallery items")
	}

	if err := s.Repository.DeleteGallery(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete gallery")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Delete Gallery"),
		Entity:    utils.NullString("Gallery"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

func (s *Service) AddGalleryItemService(galleryID uuid.UUID, req requests.AddGalleryItemRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	id, ferr := s.createGalleryItem(requests.CreateGalleryItemRequest{
		GalleryID: galleryID,
		MediaID:   req.MediaID,
		SortOrder: req.SortOrder,
	}, c)
	if ferr != nil {
		return ferr
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Add Gallery Item"),
		Entity:    utils.NullString("Gallery Item"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

func (s *Service) DeleteGalleryItemService(itemID uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	if err := s.Repository.DeleteGalleryItem(itemID, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete gallery item")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Delete Gallery Item"),
		Entity:    utils.NullString("Gallery Item"),
		EntityID:  utils.NullUUID(itemID),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}
