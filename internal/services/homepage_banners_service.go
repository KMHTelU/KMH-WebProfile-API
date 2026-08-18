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

	bannerID, ferr := s.createHomepageBanner(requests.HomepageBannerJSONRequest{
		Title:     request.Title,
		Subtitle:  request.Subtitle,
		MediaID:   media.ID,
		CtaText:   request.CtaText,
		CtaUrl:    request.CtaUrl,
		IsActive:  request.IsActive,
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	}, c)
	if ferr != nil {
		return ferr
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Create Homepage Banner", Valid: true},
		Entity:    sql.NullString{String: "HomepageBanner with ID: " + bannerID.String(), Valid: true},
		EntityID:  uuid.NullUUID{UUID: bannerID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
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

	banner, err := s.Repository.GetHomepageBannerByID(id, c)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Homepage banner not found")
	}

	// Banner harus dihapus lebih dulu sebelum media-nya karena homepage_banners
	// punya foreign key ke media; urutan sebaliknya akan melanggar constraint.
	if err := s.Repository.DeleteHomepageBanner(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete homepage banner")
	}

	if banner.MediaID.Valid {
		if err := s.DeleteMediaService(banner.MediaID.UUID, c); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete associated media")
		}
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Homepage Banner", Valid: true},
		Entity:    sql.NullString{String: "HomepageBanner with ID: " + id.String(), Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
