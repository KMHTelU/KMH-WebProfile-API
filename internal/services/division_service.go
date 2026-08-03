package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateDivisionService(req requests.CreateDivisionRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	id, ferr := s.createDivision(req, c)
	if ferr != nil {
		return ferr
	}

	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Create Division"),
		Entity:    utils.NullString("Division with Slug: " + req.Slug),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

// GetDivisionByIDService bersifat publik (read-only) untuk halaman profil.
func (s *Service) GetDivisionByIDService(id uuid.UUID, c fiber.Ctx) (generated.GetDivisionByIDRow, *fiber.Error) {
	division, err := s.Repository.GetDivisionByID(id, c)
	if err != nil {
		return generated.GetDivisionByIDRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get division")
	}
	return division, nil
}

// GetAllDivisionsService bersifat publik (read-only) untuk halaman profil.
func (s *Service) GetAllDivisionsService(c fiber.Ctx) ([]generated.GetAllDivisionsRow, *fiber.Error) {
	divisions, err := s.Repository.GetAllDivisions(c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get divisions")
	}
	return divisions, nil
}

func (s *Service) UpdateDivisionService(id uuid.UUID, req requests.UpdateDivisionRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	if ferr := s.updateDivision(id, req, c); ferr != nil {
		return ferr
	}

	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Update Division"),
		Entity:    utils.NullString("Division with Slug: " + req.Slug),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) DeleteDivisionService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	if err := s.Repository.DeleteDivision(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete division")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Delete Division"),
		Entity:    utils.NullString("Division"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) UpdateDivisionIconService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
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
		IconMediaID: utils.NullUUID(media.ID),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update division icon")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Update Division Icon"),
		Entity:    utils.NullString("Division with ID: " + id.String()),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
