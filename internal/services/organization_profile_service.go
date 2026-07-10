package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateOrganizationProfileService(req requests.CreateOrganizationProfileRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	var params generated.InsertOrganizationProfileParams = generated.InsertOrganizationProfileParams{
		ID:           uuid.New(),
		Name:         sql.NullString{String: req.Name, Valid: req.Name != ""},
		ShortName:    sql.NullString{String: req.ShortName, Valid: req.ShortName != ""},
		Description:  sql.NullString{String: req.Description, Valid: req.Description != ""},
		Vision:       sql.NullString{String: req.Vision, Valid: req.Vision != ""},
		Mission:      sql.NullString{String: req.Mission, Valid: req.Mission != ""},
		History:      sql.NullString{String: req.History, Valid: req.History != ""},
		Address:      sql.NullString{String: req.Address, Valid: req.Address != ""},
		Email:        sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		InstagramUrl: sql.NullString{String: req.InstagramUrl, Valid: req.InstagramUrl != ""},
		YoutubeUrl:   sql.NullString{String: req.YoutubeUrl, Valid: req.YoutubeUrl != ""},
		WebsiteUrl:   sql.NullString{String: req.WebsiteUrl, Valid: req.WebsiteUrl != ""},
	}

	if err := s.Repository.CreateOrganizationProfile(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create organization profile")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Create Organization Profile", Valid: true},
		Entity:    sql.NullString{String: "Organization Profile", Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) GetOrganizationProfileService(id uuid.UUID, c fiber.Ctx) (generated.GetOrganizationProfileRow, *fiber.Error) {
	profile, err := s.Repository.GetOrganizationProfile(id, c)
	if err != nil {
		return generated.GetOrganizationProfileRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get organization profile")
	}
	return profile, nil
}

func (s *Service) UpdateOrganizationProfileService(id uuid.UUID, req requests.UpdateOrganizationProfileRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	var params generated.UpdateOrganizationProfileParams = generated.UpdateOrganizationProfileParams{
		ID:           id,
		Name:         sql.NullString{String: req.Name, Valid: req.Name != ""},
		ShortName:    sql.NullString{String: req.ShortName, Valid: req.ShortName != ""},
		Description:  sql.NullString{String: req.Description, Valid: req.Description != ""},
		Vision:       sql.NullString{String: req.Vision, Valid: req.Vision != ""},
		Mission:      sql.NullString{String: req.Mission, Valid: req.Mission != ""},
		History:      sql.NullString{String: req.History, Valid: req.History != ""},
		Address:      sql.NullString{String: req.Address, Valid: req.Address != ""},
		Email:        sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		InstagramUrl: sql.NullString{String: req.InstagramUrl, Valid: req.InstagramUrl != ""},
		YoutubeUrl:   sql.NullString{String: req.YoutubeUrl, Valid: req.YoutubeUrl != ""},
		WebsiteUrl:   sql.NullString{String: req.WebsiteUrl, Valid: req.WebsiteUrl != ""},
	}
	if err := s.Repository.UpdateOrganizationProfile(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update organization profile")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Organization Profile", Valid: true},
		Entity:    sql.NullString{String: "Organization Profile", Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) DeleteOrganizationProfileService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	if err := s.Repository.DeleteOrganizationProfile(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete organization profile")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Organization Profile", Valid: true},
		Entity:    sql.NullString{String: "Organization Profile", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) UploadOrganizationProfileLogoService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	photo, err := c.FormFile("logo")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to get logo")
	}

	media, erro := s.UploadMediaService(photo, requests.CreateMediaRequest{
		FileName: photo.Filename,
		FileType: photo.Header.Get("Content-Type"),
		MimeType: photo.Header.Get("Content-Type"),
		FileSize: photo.Size,
		AltText:  "Logo of organization profile",
		Caption:  "Logo of organization profile",
	}, c)
	if erro != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload logo")
	}

	if err := s.Repository.UpdateOrganizationProfileLogo(generated.UpdateOrganizationProfileLogoParams{
		ID:          id,
		LogoMediaID: uuid.NullUUID{UUID: media.ID, Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update organization profile logo")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Organization Profile Logo", Valid: true},
		Entity:    sql.NullString{String: "Organization Profile", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
