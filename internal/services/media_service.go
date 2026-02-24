package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) UploadMediaService(file interface{}, request requests.CreateMediaRequest, c fiber.Ctx) (generated.Medium, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return generated.Medium{}, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	mediaID := uuid.New()

	result, err := s.Cloudinary.Upload.Upload(c, file, uploader.UploadParams{
		Folder:   "kmh_media",
		PublicID: "media-" + mediaID.String(),
	})
	if err != nil {
		return generated.Medium{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to upload media")
	}

	media, err := s.Repository.InsertMedia(generated.InsertMediaParams{
		ID:         mediaID,
		Url:        sql.NullString{String: result.SecureURL, Valid: true},
		FileType:   sql.NullString{String: request.FileType, Valid: true},
		FileName:   sql.NullString{String: request.FileName, Valid: true},
		MimeType:   sql.NullString{String: request.MimeType, Valid: true},
		FileSize:   sql.NullInt64{Int64: request.FileSize, Valid: true},
		AltText:    sql.NullString{String: request.AltText, Valid: true},
		Caption:    sql.NullString{String: request.Caption, Valid: true},
		UploadedBy: uuid.NullUUID{UUID: claim.UserID, Valid: true},
	}, c)
	if err != nil {
		return generated.Medium{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to save media record")
	}

	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID},
		Action:    sql.NullString{String: "Upload Media", Valid: true},
		Entity:    sql.NullString{String: "Media with ID: " + mediaID.String(), Valid: true},
		EntityID:  uuid.NullUUID{UUID: mediaID},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return generated.Medium{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}

	return media, nil
}

func (s *Service) DeleteMediaService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	var invalidate bool = true
	if _, err := s.Cloudinary.Upload.Destroy(c, uploader.DestroyParams{
		PublicID:   "media-" + id.String(),
		Invalidate: &invalidate,
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to destroy media in Cloudinary: "+err.Error())
	}
	if err := s.Repository.DeleteMedia(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete media")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Media", Valid: true},
		Entity:    sql.NullString{String: "Media with ID: " + id.String(), Valid: true},
		EntityID:  uuid.NullUUID{UUID: id},
		IpAddress: sql.NullString{String: c.IP()},
		UserAgent: sql.NullString{String: c.UserAgent()},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
