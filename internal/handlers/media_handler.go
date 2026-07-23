package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

// UploadMediaHandler adalah endpoint upload media generik.
// Dipakai admin untuk memperoleh media_id yang lalu dipasang ke
// blog-posts (featured_media_id), events (cover_media_id), atau
// gallery items (media_id) — resource yang tidak punya endpoint upload khusus.
//
// Request : multipart/form-data, field file bernama "file".
// Response: objek media { id, url, file_type, mime_type, file_name, ... }.
func (h *Handler) UploadMediaHandler(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Failed to get file")
	}

	contentType := file.Header.Get("Content-Type")
	media, ferr := h.Service.UploadMediaService(file, requests.CreateMediaRequest{
		FileName: file.Filename,
		FileType: contentType,
		MimeType: contentType,
		FileSize: file.Size,
		AltText:  file.Filename,
		Caption:  "",
	}, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}

	return utils.RespondWithCreated(c, "Media uploaded successfully", media)
}
