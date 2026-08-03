package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateMemberDivisionService(req requests.CreateMemberDivisionRequest, c fiber.Ctx) (uuid.UUID, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	id, ferr := s.createMemberDivision(req, c)
	if ferr != nil {
		return uuid.Nil, ferr
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Create Member Division"),
		Entity:    utils.NullString("Member Division"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return id, nil
}

// createMemberDivision menyimpan satu penugasan tanpa mencatat activity log,
// supaya alur bulk dan import bisa memakai ulang validasinya lalu mencatat satu
// log ringkasan saja.
func (s *Service) createMemberDivision(req requests.CreateMemberDivisionRequest, c fiber.Ctx) (uuid.UUID, *fiber.Error) {
	memberID := utils.NullUUID(req.MemberID)
	divisionID := utils.NullUUID(req.DivisionID)

	if _, err := s.Repository.GetMemberDivisionByPair(memberID, divisionID, c); err == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusConflict, "Member is already assigned to this division")
	}

	id := uuid.New()
	if err := s.Repository.CreateMemberDivision(generated.InsertMemberDivisionParams{
		ID:         id,
		MemberID:   memberID,
		DivisionID: divisionID,
		RoleTitle:  utils.NullString(req.RoleTitle),
	}, c); err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to assign member to division")
	}
	return id, nil
}

func (s *Service) UpdateMemberDivisionService(id uuid.UUID, req requests.UpdateMemberDivisionRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if ferr := s.updateMemberDivision(id, req, c); ferr != nil {
		return ferr
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Update Member Division"),
		Entity:    utils.NullString("Member Division"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

func (s *Service) updateMemberDivision(id uuid.UUID, req requests.UpdateMemberDivisionRequest, c fiber.Ctx) *fiber.Error {
	if _, err := s.Repository.GetMemberDivisionByID(id, c); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Member division not found")
	}
	if err := s.Repository.UpdateMemberDivision(generated.UpdateMemberDivisionParams{
		ID:         id,
		MemberID:   utils.NullUUID(req.MemberID),
		DivisionID: utils.NullUUID(req.DivisionID),
		RoleTitle:  utils.NullString(req.RoleTitle),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update member division")
	}
	return nil
}

func (s *Service) DeleteMemberDivisionService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteMemberDivision(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete member division")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Delete Member Division"),
		Entity:    utils.NullString("Member Division"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

// GetDivisionsByMemberService bersifat publik untuk halaman profil anggota.
func (s *Service) GetDivisionsByMemberService(memberID uuid.UUID, c fiber.Ctx) ([]generated.GetMemberDivisionsByMemberIDRow, *fiber.Error) {
	rows, err := s.Repository.ListDivisionsByMemberID(utils.NullUUID(memberID), c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get divisions of member")
	}
	return rows, nil
}

// GetMembersByDivisionService bersifat publik untuk halaman profil divisi.
func (s *Service) GetMembersByDivisionService(divisionID uuid.UUID, c fiber.Ctx) ([]generated.GetMemberDivisionsByDivisionIDRow, *fiber.Error) {
	rows, err := s.Repository.ListMembersByDivisionID(utils.NullUUID(divisionID), c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get members of division")
	}
	return rows, nil
}
