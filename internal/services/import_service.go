package services

import (
	"fmt"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/importer"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// rowHandler mengubah satu baris berkas menjadi baris basis data.
// Pada mode preview handler hanya memvalidasi tanpa menulis apa pun.
type rowHandler func(row importer.Row, dryRun bool, seen seenKeys, c fiber.Ctx) (uuid.UUID, error)

// seenKeys mencatat nilai unik yang sudah muncul pada berkas yang sedang diproses.
//
// Pemeriksaan ke basis data saja tidak cukup: dua baris dengan NIM sama dalam
// satu berkas akan lolos karena baris pertama belum tersimpan saat baris kedua
// diperiksa, dan pada mode preview tidak ada yang tersimpan sama sekali.
type seenKeys map[string]int

// claim mengembalikan pesan bila nilai sudah dipakai baris sebelumnya, lalu
// mencatat nilai tersebut.
func (s seenKeys) claim(field, value string, rowNumber int) error {
	key := field + "=" + value
	if previous, exists := s[key]; exists {
		return fmt.Errorf("%s %s sudah dipakai pada baris %d di berkas ini", field, value, previous)
	}
	s[key] = rowNumber
	return nil
}

// ImportTemplateService menghasilkan berkas template untuk sebuah entitas.
func (s *Service) ImportTemplateService(entity, format string, c fiber.Ctx) ([]byte, string, *fiber.Error) {
	if _, ferr := s.requireAuth(c); ferr != nil {
		return nil, "", ferr
	}

	spec, found := importer.SpecFor(entity)
	if !found {
		return nil, "", fiber.NewError(fiber.StatusNotFound, "Import template is not available for this resource")
	}

	switch format {
	case "", "xlsx":
		content, err := importer.BuildXLSX(spec)
		if err != nil {
			return nil, "", fiber.NewError(fiber.StatusInternalServerError, "Failed to build template")
		}
		return content, fmt.Sprintf("template_import_%s.xlsx", spec.Key), nil
	case "csv":
		content, err := importer.BuildCSV(spec)
		if err != nil {
			return nil, "", fiber.NewError(fiber.StatusInternalServerError, "Failed to build template")
		}
		return content, fmt.Sprintf("template_import_%s.csv", spec.Key), nil
	default:
		return nil, "", fiber.NewError(fiber.StatusBadRequest, "Unsupported format, use xlsx or csv")
	}
}

// ImportService membaca berkas unggahan lalu memproses tiap baris.
// Ketika dryRun bernilai true, tidak ada satu pun baris yang ditulis ke basis data.
func (s *Service) ImportService(entity, fileName string, content []byte, dryRun bool, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}

	spec, found := importer.SpecFor(entity)
	if !found {
		return nil, fiber.NewError(fiber.StatusNotFound, "Import is not available for this resource")
	}
	handler, found := s.rowHandlerFor(entity, claim)
	if !found {
		return nil, fiber.NewError(fiber.StatusNotFound, "Import is not available for this resource")
	}

	rows, err := importer.Parse(spec, fileName, content, utils.MaxBulkItems)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	report := utils.NewBulkReport(len(rows))
	seen := make(seenKeys)
	for index, row := range rows {
		id, err := handler(row, dryRun, seen, c)
		if err != nil {
			report.AddRowFailure(index, row.Number, fmt.Sprintf("baris %d: %s", row.Number, err.Error()))
			continue
		}
		report.AddRowSuccess(index, row.Number, id)
	}

	if !dryRun {
		action := fmt.Sprintf("Import %s", spec.Key)
		s.logBulk(claim, action, spec.Title, report, c)
	}

	return report, nil
}

func (s *Service) rowHandlerFor(entity string, claim *utils.Claims) (rowHandler, bool) {
	switch entity {
	case "members":
		return s.importMemberRow, true
	case "divisions":
		return s.importDivisionRow, true
	case "member-divisions":
		return s.importMemberDivisionRow, true
	case "events":
		return func(row importer.Row, dryRun bool, seen seenKeys, c fiber.Ctx) (uuid.UUID, error) {
			return s.importEventRow(row, claim.UserID, dryRun, seen, c)
		}, true
	case "roles":
		return s.importRoleRow, true
	default:
		return nil, false
	}
}

func (s *Service) importMemberRow(row importer.Row, dryRun bool, seen seenKeys, c fiber.Ctx) (uuid.UUID, error) {
	name := row.String("name")
	nim := row.String("nim")
	if name == "" {
		return uuid.Nil, fmt.Errorf("kolom name wajib diisi")
	}
	if nim == "" {
		return uuid.Nil, fmt.Errorf("kolom nim wajib diisi")
	}
	if err := seen.claim("nim", nim, row.Number); err != nil {
		return uuid.Nil, err
	}

	periodStart, err := row.Int32("period_start")
	if err != nil {
		return uuid.Nil, err
	}
	periodEnd, err := row.Int32("period_end")
	if err != nil {
		return uuid.Nil, err
	}
	if periodStart == 0 || periodEnd == 0 {
		return uuid.Nil, fmt.Errorf("kolom period_start dan period_end wajib diisi")
	}
	if periodEnd < periodStart {
		return uuid.Nil, fmt.Errorf("period_end (%d) lebih awal dari period_start (%d)", periodEnd, periodStart)
	}

	if existing, err := s.Repository.GetMemberByNIM(nim, c); err == nil {
		return uuid.Nil, fmt.Errorf("NIM %s sudah terdaftar atas nama %s", nim, existing.Name.String)
	}

	req := requests.CreateMemberRequest{
		Name:         name,
		Nim:          nim,
		Email:        row.String("email"),
		Phone:        row.String("phone"),
		Bio:          row.String("bio"),
		InstagramUrl: row.String("instagram_url"),
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
	}
	if dryRun {
		return uuid.Nil, nil
	}

	id, ferr := s.createMember(req, c)
	if ferr != nil {
		return uuid.Nil, fmt.Errorf("%s", ferr.Message)
	}
	return id, nil
}

func (s *Service) importDivisionRow(row importer.Row, dryRun bool, seen seenKeys, c fiber.Ctx) (uuid.UUID, error) {
	name := row.String("name")
	slug := row.String("slug")
	if name == "" {
		return uuid.Nil, fmt.Errorf("kolom name wajib diisi")
	}
	if slug == "" {
		return uuid.Nil, fmt.Errorf("kolom slug wajib diisi")
	}
	if err := seen.claim("slug", slug, row.Number); err != nil {
		return uuid.Nil, err
	}

	if _, err := s.Repository.GetDivisionBySlug(slug, c); err == nil {
		return uuid.Nil, fmt.Errorf("slug divisi %s sudah dipakai", slug)
	}

	var coordinatorID uuid.UUID
	if nim := row.String("coordinator_nim"); nim != "" {
		coordinator, err := s.Repository.GetMemberByNIM(nim, c)
		if err != nil {
			return uuid.Nil, fmt.Errorf("coordinator_nim %s tidak ditemukan pada data anggota", nim)
		}
		coordinatorID = coordinator.ID
	}

	req := requests.CreateDivisionRequest{
		Name:          name,
		Slug:          slug,
		Description:   row.String("description"),
		CoordinatorID: coordinatorID,
	}
	if dryRun {
		return uuid.Nil, nil
	}

	id, ferr := s.createDivision(req, c)
	if ferr != nil {
		return uuid.Nil, fmt.Errorf("%s", ferr.Message)
	}
	return id, nil
}

func (s *Service) importMemberDivisionRow(row importer.Row, dryRun bool, seen seenKeys, c fiber.Ctx) (uuid.UUID, error) {
	nim := row.String("member_nim")
	slug := row.String("division_slug")
	if nim == "" {
		return uuid.Nil, fmt.Errorf("kolom member_nim wajib diisi")
	}
	if slug == "" {
		return uuid.Nil, fmt.Errorf("kolom division_slug wajib diisi")
	}
	if err := seen.claim("pasangan member_nim dan division_slug", nim+"/"+slug, row.Number); err != nil {
		return uuid.Nil, err
	}

	member, err := s.Repository.GetMemberByNIM(nim, c)
	if err != nil {
		return uuid.Nil, fmt.Errorf("member_nim %s tidak ditemukan pada data anggota", nim)
	}
	division, err := s.Repository.GetDivisionBySlug(slug, c)
	if err != nil {
		return uuid.Nil, fmt.Errorf("division_slug %s tidak ditemukan pada data divisi", slug)
	}

	if _, err := s.Repository.GetMemberDivisionByPair(utils.NullUUID(member.ID), utils.NullUUID(division.ID), c); err == nil {
		return uuid.Nil, fmt.Errorf("anggota %s sudah terdaftar di divisi %s", nim, slug)
	}

	req := requests.CreateMemberDivisionRequest{
		MemberID:   member.ID,
		DivisionID: division.ID,
		RoleTitle:  row.String("role_title"),
	}
	if dryRun {
		return uuid.Nil, nil
	}

	id, ferr := s.createMemberDivision(req, c)
	if ferr != nil {
		return uuid.Nil, fmt.Errorf("%s", ferr.Message)
	}
	return id, nil
}

func (s *Service) importEventRow(row importer.Row, createdBy uuid.UUID, dryRun bool, seen seenKeys, c fiber.Ctx) (uuid.UUID, error) {
	title := row.String("title")
	slug := row.String("slug")
	if title == "" {
		return uuid.Nil, fmt.Errorf("kolom title wajib diisi")
	}
	if slug == "" {
		return uuid.Nil, fmt.Errorf("kolom slug wajib diisi")
	}
	if err := seen.claim("slug", slug, row.Number); err != nil {
		return uuid.Nil, err
	}

	startTime, err := row.Time("start_time")
	if err != nil {
		return uuid.Nil, err
	}
	if startTime.IsZero() {
		return uuid.Nil, fmt.Errorf("kolom start_time wajib diisi")
	}
	endTime, err := row.Time("end_time")
	if err != nil {
		return uuid.Nil, err
	}
	if !endTime.IsZero() && endTime.Before(startTime) {
		return uuid.Nil, fmt.Errorf("end_time lebih awal dari start_time")
	}

	isPublished, err := row.Bool("is_published", false)
	if err != nil {
		return uuid.Nil, err
	}

	if _, err := s.Repository.GetEventBySlug(slug, c); err == nil {
		return uuid.Nil, fmt.Errorf("slug event %s sudah dipakai", slug)
	}

	req := requests.CreateEventRequest{
		Title:           title,
		Slug:            slug,
		Description:     row.String("description"),
		EventType:       row.String("event_type"),
		StartTime:       startTime,
		EndTime:         endTime,
		Location:        row.String("location"),
		GoogleMapsUrl:   row.String("google_maps_url"),
		RegistrationUrl: row.String("registration_url"),
		Status:          row.String("status"),
		IsPublished:     isPublished,
	}
	if dryRun {
		return uuid.Nil, nil
	}

	id, ferr := s.createEvent(req, createdBy, c)
	if ferr != nil {
		return uuid.Nil, fmt.Errorf("%s", ferr.Message)
	}
	return id, nil
}

func (s *Service) importRoleRow(row importer.Row, dryRun bool, seen seenKeys, c fiber.Ctx) (uuid.UUID, error) {
	name := row.String("name")
	if name == "" {
		return uuid.Nil, fmt.Errorf("kolom name wajib diisi")
	}
	if err := seen.claim("name", name, row.Number); err != nil {
		return uuid.Nil, err
	}
	if _, err := s.Repository.GetRoleByName(name, c); err == nil {
		return uuid.Nil, fmt.Errorf("role dengan nama %s sudah ada", name)
	}

	req := requests.CreateRoleRequest{
		Name:        name,
		Description: row.String("description"),
	}
	if dryRun {
		return uuid.Nil, nil
	}

	id, ferr := s.createRole(req, c)
	if ferr != nil {
		return uuid.Nil, fmt.Errorf("%s", ferr.Message)
	}
	return id, nil
}
