package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateMemberService(req requests.CreateMemberRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	id, ferr := s.createMember(req, c)
	if ferr != nil {
		return ferr
	}

	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Create Member"),
		Entity:    utils.NullString("Member with NIM: " + req.Nim),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
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
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	if ferr := s.updateMember(id, req, c); ferr != nil {
		return ferr
	}

	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Update Member"),
		Entity:    utils.NullString("Member with NIM: " + req.Nim),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) DeleteMemberService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	// Lepaskan keterkaitan lebih dulu supaya tidak melanggar foreign key:
	// penugasan divisi milik anggota ini dan posisi koordinator yang ia pegang.
	memberRef := utils.NullUUID(id)
	if err := s.Repository.DeleteMemberDivisionsByMemberID(memberRef, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal melepas penugasan divisi milik anggota ini")
	}
	if err := s.Repository.ClearDivisionCoordinator(memberRef, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal melepas posisi koordinator milik anggota ini")
	}

	if err := s.Repository.DeleteMember(id, c); err != nil {
		if utils.IsForeignKeyViolation(err) {
			return fiber.NewError(fiber.StatusConflict,
				"Anggota tidak bisa dihapus karena masih dirujuk data lain. Lepaskan keterkaitan tersebut terlebih dahulu, lalu coba lagi.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus anggota")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Delete Member"),
		Entity:    utils.NullString("Member"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}

func (s *Service) UploadMemberPhotoService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
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
		PhotoMediaID: utils.NullUUID(media.ID),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update member photo")
	}
	if err := s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Update Member Photo"),
		Entity:    utils.NullString("Member with ID: " + id.String()),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create log")
	}
	return nil
}
