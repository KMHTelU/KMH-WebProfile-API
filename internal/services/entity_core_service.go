package services

// Fungsi di berkas ini menyimpan satu baris tanpa memeriksa token dan tanpa
// mencatat activity log. Endpoint tunggal, bulk, dan import sama-sama memakainya
// supaya pemetaan field beserta pemeriksaan duplikat hanya ditulis sekali.

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// marshalJSONList menyimpan slice sebagai JSONB. Slice nil dianggap daftar
// kosong supaya kolom tidak pernah berisi literal null.
func marshalJSONList(v any) json.RawMessage {
	raw, err := json.Marshal(v)
	if err != nil || string(raw) == "null" {
		return json.RawMessage("[]")
	}
	return raw
}

func (s *Service) createMember(req requests.CreateMemberRequest, c fiber.Ctx) (uuid.UUID, *fiber.Error) {
	if existing, err := s.Repository.GetMemberByNIM(req.Nim, c); err == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusConflict,
			fmt.Sprintf("NIM %s sudah terdaftar atas nama %s", req.Nim, existing.Name.String))
	}

	id := uuid.New()
	if err := s.Repository.CreateMember(generated.InsertMemberParams{
		ID:           id,
		Name:         utils.NullString(req.Name),
		Nim:          utils.NullString(req.Nim),
		Email:        utils.NullString(req.Email),
		Phone:        utils.NullString(req.Phone),
		Bio:          utils.NullString(req.Bio),
		InstagramUrl: utils.NullString(req.InstagramUrl),
		Faculty:      utils.NullString(req.Faculty),
		StudyProgram: utils.NullString(req.StudyProgram),
		CohortYear:   sql.NullInt32{Int32: req.CohortYear, Valid: req.CohortYear != 0},
		PeriodStart:  req.PeriodStart,
		PeriodEnd:    req.PeriodEnd,
		IsActive:     utils.NullBool(true),
	}, c); err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create member")
	}
	return id, nil
}

func (s *Service) updateMember(id uuid.UUID, req requests.UpdateMemberRequest, c fiber.Ctx) *fiber.Error {
	existing, err := s.Repository.GetMemberByID(id, c)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Member not found")
	}

	// NIM boleh tetap sama untuk anggota yang bersangkutan, tetapi tidak boleh
	// menabrak NIM anggota lain.
	if req.Nim != "" {
		if other, err := s.Repository.GetMemberByNIM(req.Nim, c); err == nil && other.ID != id {
			return fiber.NewError(fiber.StatusConflict,
				fmt.Sprintf("NIM %s sudah dipakai anggota lain", req.Nim))
		}
	}

	name := existing.Name
	if req.Name != "" {
		name = utils.NullString(req.Name)
	}
	nim := existing.Nim
	if req.Nim != "" {
		nim = utils.NullString(req.Nim)
	}
	email := existing.Email
	if req.Email != "" {
		email = utils.NullString(req.Email)
	}
	phone := existing.Phone
	if req.Phone != "" {
		phone = utils.NullString(req.Phone)
	}
	bio := existing.Bio
	if req.Bio != "" {
		bio = utils.NullString(req.Bio)
	}
	instagram := existing.InstagramUrl
	if req.InstagramUrl != "" {
		instagram = utils.NullString(req.InstagramUrl)
	}
	faculty := existing.Faculty
	if req.Faculty != "" {
		faculty = utils.NullString(req.Faculty)
	}
	studyProgram := existing.StudyProgram
	if req.StudyProgram != "" {
		studyProgram = utils.NullString(req.StudyProgram)
	}
	cohortYear := existing.CohortYear
	if req.CohortYear != 0 {
		cohortYear = sql.NullInt32{Int32: req.CohortYear, Valid: true}
	}
	periodStart := existing.PeriodStart
	if req.PeriodStart != 0 {
		periodStart = req.PeriodStart
	}
	periodEnd := existing.PeriodEnd
	if req.PeriodEnd != 0 {
		periodEnd = req.PeriodEnd
	}
	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = sql.NullBool{Bool: *req.IsActive, Valid: true}
	}

	if err := s.Repository.UpdateMember(generated.UpdateMemberParams{
		ID:           id,
		Name:         name,
		Nim:          nim,
		Email:        email,
		Phone:        phone,
		Bio:          bio,
		InstagramUrl: instagram,
		Faculty:      faculty,
		StudyProgram: studyProgram,
		CohortYear:   cohortYear,
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
		IsActive:     isActive,
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update member")
	}
	return nil
}

func (s *Service) createDivision(req requests.CreateDivisionRequest, c fiber.Ctx) (uuid.UUID, *fiber.Error) {
	if _, err := s.Repository.GetDivisionBySlug(req.Slug, c); err == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusConflict,
			fmt.Sprintf("Slug divisi %s sudah dipakai", req.Slug))
	}

	id := uuid.New()
	if err := s.Repository.CreateDivision(generated.InsertDivisionParams{
		ID:               id,
		Name:             utils.NullString(req.Name),
		Slug:             utils.NullString(req.Slug),
		Subtitle:         utils.NullString(req.Subtitle),
		Description:      utils.NullString(req.Description),
		Responsibilities: marshalJSONList(req.Responsibilities),
		Programs:         marshalJSONList(req.Programs),
		CoordinatorID:    utils.NullUUID(req.CoordinatorID),
		IsActive:         utils.NullBool(true),
	}, c); err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create division")
	}
	return id, nil
}

func (s *Service) updateDivision(id uuid.UUID, req requests.UpdateDivisionRequest, c fiber.Ctx) *fiber.Error {
	existing, err := s.Repository.GetDivisionByID(id, c)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Division not found")
	}

	if req.Slug != "" {
		if other, err := s.Repository.GetDivisionBySlug(req.Slug, c); err == nil && other.ID != id {
			return fiber.NewError(fiber.StatusConflict,
				fmt.Sprintf("Slug divisi %s sudah dipakai divisi lain", req.Slug))
		}
	}

	name := existing.Name
	if req.Name != "" {
		name = utils.NullString(req.Name)
	}
	slug := existing.Slug
	if req.Slug != "" {
		slug = utils.NullString(req.Slug)
	}
	subtitle := existing.Subtitle
	if req.Subtitle != "" {
		subtitle = utils.NullString(req.Subtitle)
	}
	description := existing.Description
	if req.Description != "" {
		description = utils.NullString(req.Description)
	}
	// Slice nil berarti field tidak dikirim; slice kosong berarti admin sengaja
	// mengosongkan daftar, jadi tetap menimpa nilai lama.
	responsibilities := existing.Responsibilities
	if req.Responsibilities != nil {
		responsibilities = marshalJSONList(req.Responsibilities)
	}
	programs := existing.Programs
	if req.Programs != nil {
		programs = marshalJSONList(req.Programs)
	}
	coordinatorID := existing.CoordinatorID
	if req.ClearCoordinator {
		// Lepas koordinator: kolom di-set NULL.
		coordinatorID = uuid.NullUUID{}
	} else if req.CoordinatorID != uuid.Nil {
		coordinatorID = utils.NullUUID(req.CoordinatorID)
	}
	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = sql.NullBool{Bool: *req.IsActive, Valid: true}
	}

	if err := s.Repository.UpdateDivision(generated.UpdateDivisionParams{
		ID:               id,
		Name:             name,
		Slug:             slug,
		Subtitle:         subtitle,
		Description:      description,
		Responsibilities: responsibilities,
		Programs:         programs,
		CoordinatorID:    coordinatorID,
		IsActive:         isActive,
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update division")
	}
	return nil
}

func (s *Service) createEvent(req requests.CreateEventRequest, createdBy uuid.UUID, c fiber.Ctx) (uuid.UUID, *fiber.Error) {
	if _, err := s.Repository.GetEventBySlug(req.Slug, c); err == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusConflict,
			fmt.Sprintf("Slug event %s sudah dipakai", req.Slug))
	}

	id := uuid.New()
	if err := s.Repository.CreateEvent(generated.InsertEventParams{
		ID:              id,
		Title:           utils.NullString(req.Title),
		Slug:            utils.NullString(req.Slug),
		Description:     utils.NullString(req.Description),
		EventType:       utils.NullString(req.EventType),
		StartTime:       utils.NullTime(req.StartTime),
		EndTime:         utils.NullTime(req.EndTime),
		Location:        utils.NullString(req.Location),
		GoogleMapsUrl:   utils.NullString(req.GoogleMapsUrl),
		RegistrationUrl: utils.NullString(req.RegistrationUrl),
		CoverMediaID:    utils.NullUUID(req.CoverMediaID),
		Status:          utils.NullString(req.Status),
		IsPublished:     utils.NullBool(req.IsPublished),
		DivisionID:      utils.NullUUID(req.DivisionID),
		CreatedBy:       utils.NullUUID(createdBy),
	}, c); err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create event")
	}
	return id, nil
}

func (s *Service) updateEvent(id uuid.UUID, req requests.UpdateEventRequest, c fiber.Ctx) *fiber.Error {
	if req.Slug != "" {
		if existing, err := s.Repository.GetEventBySlug(req.Slug, c); err == nil && existing.ID != id {
			return fiber.NewError(fiber.StatusConflict,
				fmt.Sprintf("Slug event %s sudah dipakai event lain", req.Slug))
		}
	}

	if err := s.Repository.UpdateEvent(generated.UpdateEventParams{
		ID:              id,
		Title:           utils.NullString(req.Title),
		Slug:            utils.NullString(req.Slug),
		Description:     utils.NullString(req.Description),
		EventType:       utils.NullString(req.EventType),
		StartTime:       utils.NullTime(req.StartTime),
		EndTime:         utils.NullTime(req.EndTime),
		Location:        utils.NullString(req.Location),
		GoogleMapsUrl:   utils.NullString(req.GoogleMapsUrl),
		RegistrationUrl: utils.NullString(req.RegistrationUrl),
		CoverMediaID:    utils.NullUUID(req.CoverMediaID),
		Status:          utils.NullString(req.Status),
		IsPublished:     utils.NullBool(req.IsPublished),
		DivisionID:      utils.NullUUID(req.DivisionID),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update event")
	}
	return nil
}

func (s *Service) createGallery(req requests.CreateGalleryRequest, c fiber.Ctx) (generated.Gallery, *fiber.Error) {
	gallery, err := s.Repository.CreateGallery(generated.InsertGalleryParams{
		ID:          uuid.New(),
		Title:       utils.NullString(req.Title),
		Description: utils.NullString(req.Description),
		EventID:     utils.NullUUID(req.EventID),
		IsPublic:    utils.NullBool(req.IsPublic),
	}, c)
	if err != nil {
		return generated.Gallery{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to create gallery")
	}
	return gallery, nil
}

func (s *Service) updateGallery(id uuid.UUID, req requests.UpdateGalleryRequest, c fiber.Ctx) *fiber.Error {
	if err := s.Repository.UpdateGallery(generated.UpdateGalleryParams{
		ID:          id,
		Title:       utils.NullString(req.Title),
		Description: utils.NullString(req.Description),
		EventID:     utils.NullUUID(req.EventID),
		IsPublic:    utils.NullBool(req.IsPublic),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update gallery")
	}
	return nil
}

func (s *Service) createGalleryItem(req requests.CreateGalleryItemRequest, c fiber.Ctx) (uuid.UUID, *fiber.Error) {
	id := uuid.New()
	if err := s.Repository.InsertGalleryItem(generated.InsertGalleryItemParams{
		ID:        id,
		GalleryID: utils.NullUUID(req.GalleryID),
		MediaID:   utils.NullUUID(req.MediaID),
		SortOrder: utils.NullInt32(req.SortOrder),
	}, c); err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create gallery item")
	}
	return id, nil
}

func (s *Service) updateGalleryItem(id uuid.UUID, req requests.UpdateGalleryItemRequest, c fiber.Ctx) *fiber.Error {
	if _, err := s.Repository.GetGalleryItemByID(id, c); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Gallery item not found")
	}
	if err := s.Repository.UpdateGalleryItem(generated.UpdateGalleryItemParams{
		ID:        id,
		GalleryID: utils.NullUUID(req.GalleryID),
		MediaID:   utils.NullUUID(req.MediaID),
		SortOrder: utils.NullInt32(req.SortOrder),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update gallery item")
	}
	return nil
}

func (s *Service) createHomepageBanner(req requests.HomepageBannerJSONRequest, c fiber.Ctx) (uuid.UUID, *fiber.Error) {
	id := uuid.New()
	if err := s.Repository.InsertHomepageBanner(generated.InsertHomepageBannerParams{
		ID:        id,
		Title:     utils.NullString(req.Title),
		Subtitle:  utils.NullString(req.Subtitle),
		MediaID:   utils.NullUUID(req.MediaID),
		CtaText:   utils.NullString(req.CtaText),
		CtaUrl:    utils.NullString(req.CtaUrl),
		IsActive:  utils.NullBool(req.IsActive),
		StartDate: utils.NullTime(req.StartDate),
		EndDate:   utils.NullTime(req.EndDate),
	}, c); err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create homepage banner")
	}
	return id, nil
}

func (s *Service) updateHomepageBanner(id uuid.UUID, req requests.HomepageBannerJSONRequest, c fiber.Ctx) *fiber.Error {
	if _, err := s.Repository.GetHomepageBannerByID(id, c); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Homepage banner not found")
	}
	if err := s.Repository.UpdateHomepageBanner(generated.UpdateHomepageBannerParams{
		ID:        id,
		Title:     utils.NullString(req.Title),
		Subtitle:  utils.NullString(req.Subtitle),
		MediaID:   utils.NullUUID(req.MediaID),
		CtaText:   utils.NullString(req.CtaText),
		CtaUrl:    utils.NullString(req.CtaUrl),
		IsActive:  utils.NullBool(req.IsActive),
		StartDate: utils.NullTime(req.StartDate),
		EndDate:   utils.NullTime(req.EndDate),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update homepage banner")
	}
	return nil
}

func (s *Service) createRole(req requests.CreateRoleRequest, c fiber.Ctx) (uuid.UUID, *fiber.Error) {
	if _, err := s.Repository.GetRoleByName(req.Name, c); err == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusConflict,
			fmt.Sprintf("Role dengan nama %s sudah ada", req.Name))
	}

	id := uuid.New()
	if _, err := s.Repository.CreateRole(generated.InsertRoleParams{
		ID:          id,
		Name:        utils.NullString(req.Name),
		Description: utils.NullString(req.Description),
	}, c); err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create role")
	}
	return id, nil
}

func (s *Service) updateRole(id uuid.UUID, req requests.UpdateRoleRequest, c fiber.Ctx) *fiber.Error {
	if req.Name != "" {
		if existing, err := s.Repository.GetRoleByName(req.Name, c); err == nil && existing.ID != id {
			return fiber.NewError(fiber.StatusConflict,
				fmt.Sprintf("Role dengan nama %s sudah ada", req.Name))
		}
	}

	if err := s.Repository.UpdateRole(generated.UpdateRoleParams{
		ID:          id,
		Name:        utils.NullString(req.Name),
		Description: utils.NullString(req.Description),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update role")
	}
	return nil
}
