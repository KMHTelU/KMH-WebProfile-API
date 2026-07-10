package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateGalleryService(req requests.CreateGalleryRequest, c fiber.Ctx) (generated.Gallery, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return generated.Gallery{}, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.InsertGalleryParams{
		ID:          uuid.New(),
		Title:       sql.NullString{String: req.Title, Valid: req.Title != ""},
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		EventID:     uuid.NullUUID{UUID: req.EventID, Valid: req.EventID != uuid.Nil},
		IsPublic:    sql.NullBool{Bool: req.IsPublic, Valid: true},
	}

	gallery, err1 := s.Repository.CreateGallery(params, c)
	if err1 != nil {
		return generated.Gallery{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to create gallery")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Create Gallery", Valid: true},
		Entity:    sql.NullString{String: "Gallery", Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
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

func (s *Service) GetPaginatedGalleriesService(limit, offset int32, c fiber.Ctx) ([]generated.SelectAllGalleriesRow, *fiber.Error) {
	galleries, err := s.Repository.ListGalleries(generated.SelectAllGalleriesParams{
		Limit:  limit,
		Offset: offset,
	}, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get galleries")
	}
	return galleries, nil
}

func (s *Service) UpdateGalleryService(id uuid.UUID, req requests.UpdateGalleryRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.UpdateGalleryParams{
		ID:          id,
		Title:       sql.NullString{String: req.Title, Valid: req.Title != ""},
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		EventID:     uuid.NullUUID{UUID: req.EventID, Valid: req.EventID != uuid.Nil},
		IsPublic:    sql.NullBool{Bool: req.IsPublic, Valid: true},
	}

	if err := s.Repository.UpdateGallery(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update gallery")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Gallery", Valid: true},
		Entity:    sql.NullString{String: "Gallery", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) DeleteGalleryService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteGallery(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete gallery")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Gallery", Valid: true},
		Entity:    sql.NullString{String: "Gallery", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) AddGalleryItemService(galleryID uuid.UUID, req requests.AddGalleryItemRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.InsertGalleryItemParams{
		ID:        uuid.New(),
		GalleryID: uuid.NullUUID{UUID: galleryID, Valid: true},
		MediaID:   uuid.NullUUID{UUID: req.MediaID, Valid: true},
		SortOrder: sql.NullInt32{Int32: req.SortOrder, Valid: true},
	}

	if err := s.Repository.InsertGalleryItem(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to add gallery item")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Add Gallery Item", Valid: true},
		Entity:    sql.NullString{String: "Gallery Item", Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) DeleteGalleryItemService(itemID uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteGalleryItem(itemID, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete gallery item")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Gallery Item", Valid: true},
		Entity:    sql.NullString{String: "Gallery Item", Valid: true},
		EntityID:  uuid.NullUUID{UUID: itemID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}
