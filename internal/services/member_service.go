package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateMemberService(req requests.CreateMemberRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	var params generated.InsertMemberParams = generated.InsertMemberParams{
		ID:           uuid.New(),
		Name:         sql.NullString{String: req.Name, Valid: req.Name != ""},
		Npm:          sql.NullString{String: req.Npm, Valid: req.Npm != ""},
		Email:        sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		InstagramUrl: sql.NullString{String: req.InstagramUrl, Valid: req.InstagramUrl != ""},
		PeriodStart:  req.PeriodStart,
		PeriodEnd:    req.PeriodEnd,
		Bio:          sql.NullString{String: req.Bio, Valid: req.Bio != ""},
	}

	if err := s.Repository.CreateMember(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create member")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Create Member", Valid: true},
		Entity:    sql.NullString{String: "Member with NPM: " + req.Npm, Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

// GetMemberByIDService bersifat publik (read-only) untuk halaman profil.
func (s *Service) GetMemberByIDService(id uuid.UUID, c fiber.Ctx) (generated.GetMemberByIDRow, *fiber.Error) {
	member, err := s.Repository.GetMemberByID(id, c)
	if err != nil {
		return generated.GetMemberByIDRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get member")
	}
	return member, nil
}

// GetPaginatedAllMembersService bersifat publik (read-only) untuk halaman profil.
func (s *Service) GetPaginatedAllMembersService(limit, offset int32, c fiber.Ctx) ([]generated.GetAllMembersRow, *fiber.Error) {
	memberParam := generated.GetAllMembersParams{
		Limit:  limit,
		Offset: offset,
	}
	members, err := s.Repository.GetAllMembers(memberParam, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get members")
	}
	return members, nil
}

func (s *Service) UpdateMemberService(id uuid.UUID, req requests.UpdateMemberRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	var params generated.UpdateMemberParams = generated.UpdateMemberParams{
		ID:           id,
		Name:         sql.NullString{String: req.Name, Valid: req.Name != ""},
		Email:        sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		InstagramUrl: sql.NullString{String: req.InstagramUrl, Valid: req.InstagramUrl != ""},
		PeriodStart:  req.PeriodStart,
		PeriodEnd:    req.PeriodEnd,
		Bio:          sql.NullString{String: req.Bio, Valid: req.Bio != ""},
		IsActive:     sql.NullBool{Bool: req.IsActive, Valid: true},
	}
	if err := s.Repository.UpdateMember(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update member")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Member", Valid: true},
		Entity:    sql.NullString{String: "Member with NPM: " + req.Npm, Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) DeleteMemberService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	if err := s.Repository.DeleteMember(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete member")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Member", Valid: true},
		Entity:    sql.NullString{String: "Member", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) UploadMemberPhotoService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	photo, err := c.FormFile("photo")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to get photo")
	}

	media, erro := s.UploadMediaService(photo, requests.CreateMediaRequest{
		FileName: photo.Filename,
		FileType: photo.Header.Get("Content-Type"),
		MimeType: photo.Header.Get("Content-Type"),
		FileSize: photo.Size,
		AltText:  "Photo of member with ID: " + id.String(),
		Caption:  "Photo of member with ID: " + id.String(),
	}, c)
	if erro != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload photo")
	}

	if err := s.Repository.UpdateMemberPhoto(generated.UpdateMemberPhotoParams{
		ID:           id,
		PhotoMediaID: uuid.NullUUID{UUID: media.ID, Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update member photo")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Member Photo", Valid: true},
		Entity:    sql.NullString{String: "Member with ID: " + id.String(), Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
