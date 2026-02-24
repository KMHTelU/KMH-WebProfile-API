package services

import (
	"database/sql"
	"time"

	// "github.com/KMHTelU/KMH-WebProfile-API/internal/entities"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateHomepageBannerService(request requests.HomepageBannerRequest, c fiber.Ctx) *fiber.Error {
	// Implementation for creating a homepage banner
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	file, err := c.FormFile("media")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Media file is required")
	}

	media, erro := s.UploadMediaService(file, requests.CreateMediaRequest{
		FileType: file.Header.Get("Content-Type"),
		FileName: file.Filename,
		MimeType: file.Header.Get("Content-Type"),
		FileSize: file.Size,
		AltText:  request.AltText,
		Caption:  request.Caption,
	}, c)
	if erro != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload media for homepage banner")
	}

	bannerID := uuid.New()
	params := generated.InsertHomepageBannerParams{
		ID:        bannerID,
		Title:     sql.NullString{String: request.Title, Valid: true},
		Subtitle:  sql.NullString{String: request.Subtitle, Valid: true},
		MediaID:   uuid.NullUUID{UUID: media.ID, Valid: true},
		CtaText:   sql.NullString{String: request.CtaText, Valid: true},
		CtaUrl:    sql.NullString{String: request.CtaUrl, Valid: true},
		IsActive:  sql.NullBool{Bool: request.IsActive, Valid: true},
		StartDate: sql.NullTime{Time: request.StartDate, Valid: true},
		EndDate:   sql.NullTime{Time: request.EndDate, Valid: true},
	}
	if err := s.Repository.InsertHomepageBanner(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create homepage banner")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Create Homepage Banner", Valid: true},
		Entity:    sql.NullString{String: "HomepageBanner with ID: " + bannerID.String(), Valid: true},
		EntityID:  uuid.NullUUID{UUID: bannerID},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) GetHomepageBannersService(c fiber.Ctx) ([]generated.SelectAllHomepageBannersRow, *fiber.Error) {
	// Implementation for retrieving homepage banners
	rows, err := s.Repository.GetHomepageBanners(generated.SelectAllHomepageBannersParams{
		Limit:  999,
		Offset: 0,
	}, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve homepage banners")
	}
	onlyActive := make([]generated.SelectAllHomepageBannersRow, 0)
	for _, row := range rows {
		if row.IsActive.Valid && row.IsActive.Bool && (row.StartDate.Valid && row.StartDate.Time.Before(time.Now())) && (row.EndDate.Valid && row.EndDate.Time.After(time.Now())) {
			onlyActive = append(onlyActive, row)
		}
	}
	return onlyActive, nil
}

func (s *Service) GetPaginatedHomepageBannersService(limit, offset int32, c fiber.Ctx) ([]generated.SelectAllHomepageBannersRow, *fiber.Error) {
	// Implementation for retrieving homepage banners
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	rows, err := s.Repository.GetHomepageBanners(generated.SelectAllHomepageBannersParams{
		Limit:  limit,
		Offset: offset,
	}, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve homepage banners")
	}

	return rows, nil
}

func (s *Service) DeleteHomepageBannerService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	// Implementation for deleting a homepage banner
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	rows, err := s.Repository.GetHomepageBanners(generated.SelectAllHomepageBannersParams{
		Limit:  1,
		Offset: 0,
	}, c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve homepage banners")
	}

	row := rows[0]
	if row.MediaID.Valid {
		if err := s.DeleteMediaService(row.MediaID.UUID, c); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete associated media")
		}
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Delete Homepage Banner", Valid: true},
		Entity:    sql.NullString{String: "HomepageBanner with ID: " + id.String(), Valid: true},
		EntityID:  uuid.NullUUID{UUID: id},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
