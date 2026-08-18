package services

// Layanan Hall of Fame: satu endpoint publik yang mengembalikan seluruh arsip
// (generasi, orang, prestasi, timeline) untuk museum 3D, plus CRUD per entitas
// untuk admin. Foto orang memakai media yang sudah ada (photo_media_id).

import (
	"encoding/json"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type HofGenerationDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	YearStart   int32    `json:"year_start"`
	YearEnd     int32    `json:"year_end"`
	Description string   `json:"description"`
	Milestones  []string `json:"milestones"`
	Accent      string   `json:"accent"`
	SortOrder   int32    `json:"sort_order"`
}

type HofPersonDTO struct {
	ID            string   `json:"id"`
	GenerationID  string   `json:"generation_id"`
	Name          string   `json:"name"`
	Role          string   `json:"role"`
	StudyProgram  string   `json:"study_program"`
	Biography     string   `json:"biography"`
	Contributions string   `json:"contributions"`
	Legacy        string   `json:"legacy"`
	Quote         string   `json:"quote"`
	Fields        []string `json:"fields"`
	PhotoMediaID  string   `json:"photo_media_id"`
	PhotoUrl      string   `json:"photo_url"`
	SortOrder     int32    `json:"sort_order"`
}

type HofAchievementDTO struct {
	ID           string `json:"id"`
	PersonID     string `json:"person_id"`
	Title        string `json:"title"`
	Category     string `json:"category"`
	Year         int32  `json:"year"`
	Organization string `json:"organization"`
	Result       string `json:"result"`
	Description  string `json:"description"`
}

type HofTimelineDTO struct {
	ID          string `json:"id"`
	Year        string `json:"year"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sort_order"`
}

type HallOfFameResponse struct {
	Generations  []HofGenerationDTO  `json:"generations"`
	People       []HofPersonDTO      `json:"people"`
	Achievements []HofAchievementDTO `json:"achievements"`
	Timeline     []HofTimelineDTO    `json:"timeline"`
}

func decodeStringList(raw json.RawMessage) []string {
	var list []string
	if err := json.Unmarshal(raw, &list); err != nil {
		return []string{}
	}
	return list
}

// GetHallOfFameService publik — dikonsumsi museum 3D dalam satu request.
func (s *Service) GetHallOfFameService(c fiber.Ctx) (HallOfFameResponse, *fiber.Error) {
	resp := HallOfFameResponse{
		Generations:  make([]HofGenerationDTO, 0),
		People:       make([]HofPersonDTO, 0),
		Achievements: make([]HofAchievementDTO, 0),
		Timeline:     make([]HofTimelineDTO, 0),
	}

	generations, err := s.Repository.ListHofGenerations(c)
	if err != nil {
		return resp, fiber.NewError(fiber.StatusInternalServerError, "Failed to get hall of fame generations")
	}
	people, err := s.Repository.ListHofPeople(c)
	if err != nil {
		return resp, fiber.NewError(fiber.StatusInternalServerError, "Failed to get hall of fame people")
	}
	achievements, err := s.Repository.ListHofAchievements(c)
	if err != nil {
		return resp, fiber.NewError(fiber.StatusInternalServerError, "Failed to get hall of fame achievements")
	}
	timeline, err := s.Repository.ListHofTimelineEvents(c)
	if err != nil {
		return resp, fiber.NewError(fiber.StatusInternalServerError, "Failed to get hall of fame timeline")
	}

	for _, g := range generations {
		resp.Generations = append(resp.Generations, HofGenerationDTO{
			ID:          g.ID.String(),
			Name:        g.Name,
			YearStart:   g.YearStart,
			YearEnd:     g.YearEnd,
			Description: g.Description,
			Milestones:  decodeStringList(g.Milestones),
			Accent:      g.Accent,
			SortOrder:   g.SortOrder,
		})
	}
	for _, p := range people {
		photoMediaID := ""
		if p.PhotoMediaID.Valid {
			photoMediaID = p.PhotoMediaID.UUID.String()
		}
		resp.People = append(resp.People, HofPersonDTO{
			ID:            p.ID.String(),
			GenerationID:  p.GenerationID.String(),
			Name:          p.Name,
			Role:          p.Role,
			StudyProgram:  p.StudyProgram,
			Biography:     p.Biography,
			Contributions: p.Contributions,
			Legacy:        p.Legacy,
			Quote:         p.Quote,
			Fields:        decodeStringList(p.Fields),
			PhotoMediaID:  photoMediaID,
			PhotoUrl:      p.PhotoUrl.String,
			SortOrder:     p.SortOrder,
		})
	}
	for _, a := range achievements {
		resp.Achievements = append(resp.Achievements, HofAchievementDTO{
			ID:           a.ID.String(),
			PersonID:     a.PersonID.String(),
			Title:        a.Title,
			Category:     a.Category,
			Year:         a.Year,
			Organization: a.Organization,
			Result:       a.Result,
			Description:  a.Description,
		})
	}
	for _, t := range timeline {
		resp.Timeline = append(resp.Timeline, HofTimelineDTO{
			ID:          t.ID.String(),
			Year:        t.YearLabel,
			Title:       t.Title,
			Description: t.Description,
			SortOrder:   t.SortOrder,
		})
	}

	return resp, nil
}

// logHof mencatat activity log dengan pola yang sama seperti fitur lain.
func (s *Service) logHof(action, entity string, entityID uuid.UUID, userID uuid.UUID, c fiber.Ctx) {
	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(userID),
		Action:    utils.NullString(action),
		Entity:    utils.NullString(entity),
		EntityID:  utils.NullUUID(entityID),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)
}

// ── Generations ──

func (s *Service) CreateHofGenerationService(req requests.HofGenerationRequest, c fiber.Ctx) (generated.HofGeneration, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return generated.HofGeneration{}, ferr
	}
	row, err := s.Repository.CreateHofGeneration(generated.InsertHofGenerationParams{
		ID:          uuid.New(),
		Name:        req.Name,
		YearStart:   req.YearStart,
		YearEnd:     req.YearEnd,
		Description: req.Description,
		Milestones:  marshalJSONList(req.Milestones),
		Accent:      req.Accent,
		SortOrder:   req.SortOrder,
	}, c)
	if err != nil {
		return generated.HofGeneration{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to create hof generation")
	}
	s.logHof("Create HOF Generation", "HofGeneration", row.ID, claim.UserID, c)
	return row, nil
}

func (s *Service) UpdateHofGenerationService(id uuid.UUID, req requests.HofGenerationRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	if err := s.Repository.UpdateHofGeneration(generated.UpdateHofGenerationParams{
		ID:          id,
		Name:        req.Name,
		YearStart:   req.YearStart,
		YearEnd:     req.YearEnd,
		Description: req.Description,
		Milestones:  marshalJSONList(req.Milestones),
		Accent:      req.Accent,
		SortOrder:   req.SortOrder,
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update hof generation")
	}
	s.logHof("Update HOF Generation", "HofGeneration", id, claim.UserID, c)
	return nil
}

func (s *Service) DeleteHofGenerationService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	// ON DELETE CASCADE menghapus orang + prestasi di dalam generasi ini.
	if err := s.Repository.DeleteHofGeneration(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete hof generation")
	}
	s.logHof("Delete HOF Generation", "HofGeneration", id, claim.UserID, c)
	return nil
}

// ── People ──

func (s *Service) CreateHofPersonService(req requests.HofPersonRequest, c fiber.Ctx) (generated.HofPerson, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return generated.HofPerson{}, ferr
	}
	row, err := s.Repository.CreateHofPerson(generated.InsertHofPersonParams{
		ID:            uuid.New(),
		GenerationID:  req.GenerationID,
		Name:          req.Name,
		Role:          req.Role,
		StudyProgram:  req.StudyProgram,
		Biography:     req.Biography,
		Contributions: req.Contributions,
		Legacy:        req.Legacy,
		Quote:         req.Quote,
		Fields:        marshalJSONList(req.Fields),
		PhotoMediaID:  utils.NullUUID(req.PhotoMediaID),
		SortOrder:     req.SortOrder,
	}, c)
	if err != nil {
		return generated.HofPerson{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to create hof person")
	}
	s.logHof("Create HOF Person", "HofPerson", row.ID, claim.UserID, c)
	return row, nil
}

func (s *Service) UpdateHofPersonService(id uuid.UUID, req requests.HofPersonRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	if err := s.Repository.UpdateHofPerson(generated.UpdateHofPersonParams{
		ID:            id,
		GenerationID:  req.GenerationID,
		Name:          req.Name,
		Role:          req.Role,
		StudyProgram:  req.StudyProgram,
		Biography:     req.Biography,
		Contributions: req.Contributions,
		Legacy:        req.Legacy,
		Quote:         req.Quote,
		Fields:        marshalJSONList(req.Fields),
		PhotoMediaID:  utils.NullUUID(req.PhotoMediaID),
		SortOrder:     req.SortOrder,
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update hof person")
	}
	s.logHof("Update HOF Person", "HofPerson", id, claim.UserID, c)
	return nil
}

func (s *Service) DeleteHofPersonService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	if err := s.Repository.DeleteHofPerson(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete hof person")
	}
	s.logHof("Delete HOF Person", "HofPerson", id, claim.UserID, c)
	return nil
}

// ── Achievements ──

func (s *Service) CreateHofAchievementService(req requests.HofAchievementRequest, c fiber.Ctx) (generated.HofAchievement, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return generated.HofAchievement{}, ferr
	}
	category := req.Category
	if category == "" {
		category = "Other"
	}
	row, err := s.Repository.CreateHofAchievement(generated.InsertHofAchievementParams{
		ID:           uuid.New(),
		PersonID:     req.PersonID,
		Title:        req.Title,
		Category:     category,
		Year:         req.Year,
		Organization: req.Organization,
		Result:       req.Result,
		Description:  req.Description,
	}, c)
	if err != nil {
		return generated.HofAchievement{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to create hof achievement")
	}
	s.logHof("Create HOF Achievement", "HofAchievement", row.ID, claim.UserID, c)
	return row, nil
}

func (s *Service) UpdateHofAchievementService(id uuid.UUID, req requests.HofAchievementRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	category := req.Category
	if category == "" {
		category = "Other"
	}
	if err := s.Repository.UpdateHofAchievement(generated.UpdateHofAchievementParams{
		ID:           id,
		PersonID:     req.PersonID,
		Title:        req.Title,
		Category:     category,
		Year:         req.Year,
		Organization: req.Organization,
		Result:       req.Result,
		Description:  req.Description,
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update hof achievement")
	}
	s.logHof("Update HOF Achievement", "HofAchievement", id, claim.UserID, c)
	return nil
}

func (s *Service) DeleteHofAchievementService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	if err := s.Repository.DeleteHofAchievement(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete hof achievement")
	}
	s.logHof("Delete HOF Achievement", "HofAchievement", id, claim.UserID, c)
	return nil
}

// ── Timeline ──

func (s *Service) CreateHofTimelineEventService(req requests.HofTimelineEventRequest, c fiber.Ctx) (generated.HofTimelineEvent, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return generated.HofTimelineEvent{}, ferr
	}
	row, err := s.Repository.CreateHofTimelineEvent(generated.InsertHofTimelineEventParams{
		ID:          uuid.New(),
		YearLabel:   req.YearLabel,
		Title:       req.Title,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}, c)
	if err != nil {
		return generated.HofTimelineEvent{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to create hof timeline event")
	}
	s.logHof("Create HOF Timeline Event", "HofTimelineEvent", row.ID, claim.UserID, c)
	return row, nil
}

func (s *Service) UpdateHofTimelineEventService(id uuid.UUID, req requests.HofTimelineEventRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	if err := s.Repository.UpdateHofTimelineEvent(generated.UpdateHofTimelineEventParams{
		ID:          id,
		YearLabel:   req.YearLabel,
		Title:       req.Title,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update hof timeline event")
	}
	s.logHof("Update HOF Timeline Event", "HofTimelineEvent", id, claim.UserID, c)
	return nil
}

func (s *Service) DeleteHofTimelineEventService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}
	if err := s.Repository.DeleteHofTimelineEvent(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete hof timeline event")
	}
	s.logHof("Delete HOF Timeline Event", "HofTimelineEvent", id, claim.UserID, c)
	return nil
}
