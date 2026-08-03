package handlers

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

// maxImportFileSize membatasi ukuran berkas unggahan. Berkas spreadsheet berisi
// beberapa ratus baris teks jauh di bawah angka ini, jadi berkas yang lebih
// besar hampir pasti salah unggah.
const maxImportFileSize = 5 << 20 // 5 MB

func (h *Handler) DownloadImportTemplateHandler(c fiber.Ctx) error {
	entity := c.Params("entity")
	format := strings.ToLower(c.Query("format", "xlsx"))

	content, fileName, err := h.Service.ImportTemplateService(entity, format, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}

	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if format == "csv" {
		contentType = "text/csv; charset=utf-8"
	}
	c.Set(fiber.HeaderContentType, contentType)
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", fileName))
	return c.Send(content)
}

// PreviewImportHandler memvalidasi berkas tanpa menyimpan apa pun.
func (h *Handler) PreviewImportHandler(c fiber.Ctx) error {
	return h.handleImport(c, true)
}

// ImportHandler memvalidasi sekaligus menyimpan baris yang lolos.
func (h *Handler) ImportHandler(c fiber.Ctx) error {
	return h.handleImport(c, false)
}

func (h *Handler) handleImport(c fiber.Ctx, dryRun bool) error {
	entity := c.Params("entity")

	header, err := c.FormFile("file")
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "File is required, upload it as multipart form field 'file'")
	}
	if header.Size > maxImportFileSize {
		return utils.RespondWithError(c, fiber.StatusRequestEntityTooLarge, "File is too large, maximum size is 5 MB")
	}

	extension := strings.ToLower(filepath.Ext(header.Filename))
	if extension != ".xlsx" && extension != ".csv" {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Unsupported file type, use .xlsx or .csv")
	}

	file, err := header.Open()
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to read uploaded file")
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, maxImportFileSize))
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to read uploaded file")
	}

	report, ferr := h.Service.ImportService(entity, header.Filename, content, dryRun, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}

	if dryRun {
		// Preview tidak mengubah apa pun, jadi selalu 200 selama berkas terbaca.
		return utils.RespondWithOK(c, "Import preview completed, no data was saved", report)
	}
	return utils.RespondWithBulkReport(c, "Import processed", report, fiber.StatusCreated)
}
