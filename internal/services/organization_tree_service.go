package services

// Layanan struktur organisasi publik. Menggabungkan divisi + koordinatornya
// dengan jajaran pengurus inti menjadi satu respons untuk komponen
// OrganizationTree di frontend. Foto memakai foto member yang sudah ada.

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// isCoreDivision menentukan apakah sebuah divisi adalah "pengurus inti" (BPH):
// anggotanya mengisi bagian atas tree (Ketua, Wakil, Sekretaris, Bendahara)
// dan divisinya sendiri tidak ditampilkan sebagai kotak divisi.
//
// Penanda utama: division_type = 'core'. Heuristik slug/nama dipertahankan
// sebagai cadangan untuk data lama yang belum mengisi tipenya.
func isCoreDivision(slug, name, divisionType string) bool {
	if strings.EqualFold(strings.TrimSpace(divisionType), "core") {
		return true
	}
	normalize := func(v string) string {
		v = strings.ToLower(strings.TrimSpace(v))
		v = strings.ReplaceAll(v, "-", " ")
		v = strings.ReplaceAll(v, "_", " ")
		return strings.Join(strings.Fields(v), " ")
	}
	s := normalize(slug)
	n := normalize(name)

	if s == "inti" || n == "inti" || s == "bph" || n == "bph" {
		return true
	}
	for _, keyword := range []string{"pengurus inti", "pengurus harian", "badan pengurus harian"} {
		if strings.Contains(s, keyword) || strings.Contains(n, keyword) {
			return true
		}
	}
	return false
}

type OrgTreePerson struct {
	MemberID     string `json:"member_id"`
	Name         string `json:"name"`
	Nim          string `json:"nim"`
	Faculty      string `json:"faculty"`
	StudyProgram string `json:"study_program"`
	CohortYear   int32  `json:"cohort_year"`
	PhotoUrl     string `json:"photo_url"`
	RoleTitle    string `json:"role_title"`
}

type OrgTreeDivision struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Slug             string         `json:"slug"`
	Subtitle         string         `json:"subtitle"`
	Description      string         `json:"description"`
	DivisionType     string         `json:"division_type"`
	Responsibilities []string       `json:"responsibilities"`
	Coordinator      *OrgTreePerson `json:"coordinator"`
}

type OrgTreeResponse struct {
	Leadership []OrgTreePerson   `json:"leadership"`
	Divisions  []OrgTreeDivision `json:"divisions"`
}

func (s *Service) GetOrganizationTreeService(c fiber.Ctx) (OrgTreeResponse, *fiber.Error) {
	divisionRows, err := s.Repository.ListOrgTreeDivisions(c)
	if err != nil {
		return OrgTreeResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get organization tree divisions")
	}
	assignmentRows, err := s.Repository.ListOrgTreeAssignments(c)
	if err != nil {
		return OrgTreeResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get organization tree assignments")
	}

	resp := OrgTreeResponse{
		Leadership: make([]OrgTreePerson, 0),
		Divisions:  make([]OrgTreeDivision, 0),
	}

	// Penugasan dipisah: divisi inti mengisi leadership, sisanya menjadi
	// kandidat koordinator cadangan per divisi (bila coordinator_id kosong).
	fallbackCoordinator := make(map[string]OrgTreePerson)
	for _, row := range assignmentRows {
		person := OrgTreePerson{
			MemberID:     row.MemberID.String(),
			Name:         row.MemberName.String,
			Nim:          row.MemberNim.String,
			Faculty:      row.MemberFaculty.String,
			StudyProgram: row.MemberStudyProgram.String,
			CohortYear:   row.MemberCohortYear.Int32,
			PhotoUrl:     row.PhotoUrl.String,
			RoleTitle:    row.RoleTitle.String,
		}
		if isCoreDivision(row.DivisionSlug.String, row.DivisionName.String, row.DivisionType) {
			resp.Leadership = append(resp.Leadership, person)
			continue
		}
		role := strings.ToLower(row.RoleTitle.String)
		if strings.Contains(role, "koordinator") || strings.Contains(role, "ketua") {
			key := row.DivisionID.String()
			if _, exists := fallbackCoordinator[key]; !exists {
				fallbackCoordinator[key] = person
			}
		}
	}

	for _, row := range divisionRows {
		if isCoreDivision(row.Slug.String, row.Name.String, row.DivisionType) {
			continue
		}

		var responsibilities []string
		if err := json.Unmarshal(row.Responsibilities, &responsibilities); err != nil {
			responsibilities = []string{}
		}

		division := OrgTreeDivision{
			ID:               row.ID.String(),
			Name:             row.Name.String,
			Slug:             row.Slug.String,
			Subtitle:         row.Subtitle.String,
			Description:      row.Description.String,
			DivisionType:     row.DivisionType,
			Responsibilities: responsibilities,
		}
		if row.CoordinatorID.Valid {
			division.Coordinator = &OrgTreePerson{
				MemberID:     row.CoordinatorID.UUID.String(),
				Name:         row.CoordinatorName.String,
				Nim:          row.CoordinatorNim.String,
				Faculty:      row.CoordinatorFaculty.String,
				StudyProgram: row.CoordinatorStudyProgram.String,
				CohortYear:   row.CoordinatorCohortYear.Int32,
				PhotoUrl:     row.CoordinatorPhotoUrl.String,
				RoleTitle:    "Koordinator",
			}
		} else if person, ok := fallbackCoordinator[row.ID.String()]; ok {
			division.Coordinator = &person
		}

		resp.Divisions = append(resp.Divisions, division)
	}

	return resp, nil
}
