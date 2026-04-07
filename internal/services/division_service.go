package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateDivisionService(req requests.CreateDivisionRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	newId := uuid.New()
	var params generated.InsertDivisionParams = generated.InsertDivisionParams{
		ID:            newId,
		Name:          sql.NullString{String: req.Name, Valid: req.Name != ""},
		Slug:          sql.NullString{String: req.Slug, Valid: req.Slug != ""},
		Description:   sql.NullString{String: req.Description, Valid: req.Description != ""},
		CoordinatorID: uuid.NullUUID{UUID: req.CoordinatorID, Valid: req.CoordinatorID != uuid.Nil},
	}
	err = s.Repository.CreateDivision(params, c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create division")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Create Division"},
		Entity:    sql.NullString{String: "Division with Slug: " + req.Slug},
		EntityID:  uuid.NullUUID{UUID: newId},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) GetDivisionByIDService(id uuid.UUID, c fiber.Ctx) (generated.GetDivisionByIDRow, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return generated.GetDivisionByIDRow{}, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	division, err := s.Repository.GetDivisionByID(id, c)
	if err != nil {
		return generated.GetDivisionByIDRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get division")
	}
	return division, nil
}

func (s *Service) GetAllDivisionsService(c fiber.Ctx) ([]generated.GetAllDivisionsRow, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	divisions, err := s.Repository.GetAllDivisions(c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get divisions")
	}
	return divisions, nil
}

func (s *Service) UpdateDivisionService(id uuid.UUID, req requests.UpdateDivisionRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	var params generated.UpdateDivisionParams = generated.UpdateDivisionParams{
		ID:            id,
		Name:          sql.NullString{String: req.Name, Valid: req.Name != ""},
		Slug:          sql.NullString{String: req.Slug, Valid: req.Slug != ""},
		Description:   sql.NullString{String: req.Description, Valid: req.Description != ""},
		CoordinatorID: uuid.NullUUID{UUID: req.CoordinatorID, Valid: req.CoordinatorID != uuid.Nil},
		IsActive:      sql.NullBool{Bool: req.IsActive, Valid: true},
	}
	if err := s.Repository.UpdateDivision(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update division")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Update Division"},
		Entity:    sql.NullString{String: "Division with Slug: " + req.Slug},
		EntityID:  uuid.NullUUID{UUID: id},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) DeleteDivisionService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	if err := s.Repository.DeleteDivision(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete division")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Delete Division"},
		Entity:    sql.NullString{String: "Division"},
		EntityID:  uuid.NullUUID{UUID: id},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) UpdateDivisionIconService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	icon, err := c.FormFile("icon")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Icon file is required")
	}

	media, erro := s.UploadMediaService(icon, requests.CreateMediaRequest{
		FileName: icon.Filename,
		FileType: icon.Header.Get("Content-Type"),
		MimeType: icon.Header.Get("Content-Type"),
		FileSize: icon.Size,
		AltText:  "Icon division - " + id.String(),
		Caption:  "Icon division - " + id.String(),
	}, c)
	if erro != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload icon")
	}
	if err := s.Repository.UpdateDivisionIcon(generated.UpdateDivisionIconParams{
		ID:          id,
		IconMediaID: uuid.NullUUID{UUID: media.ID, Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update division icon")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Update Division Icon"},
		Entity:    sql.NullString{String: "Division with ID: " + id.String()},
		EntityID:  uuid.NullUUID{UUID: id},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
