package services

// Layanan struktur organisasi publik. Menggabungkan divisi + koordinatornya
// dengan jajaran pengurus inti menjadi satu respons untuk komponen
// OrganizationTree di frontend. Foto memakai foto member yang sudah ada.

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Divisi dengan slug berikut dianggap "pengurus inti" (BPH): anggotanya mengisi
// bagian atas tree (Ketua, Wakil, Sekretaris, Bendahara) dan divisinya sendiri
// tidak ditampilkan sebagai kotak divisi.
var orgTreeCoreSlugs = map[string]bool{
	"inti":            true,
	"pengurus-inti":   true,
	"bph":             true,
	"pengurus-harian": true,
}

type OrgTreePerson struct {
	MemberID  string `json:"member_id"`
	Name      string `json:"name"`
	Nim       string `json:"nim"`
	PhotoUrl  string `json:"photo_url"`
	RoleTitle string `json:"role_title"`
}

type OrgTreeDivision struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Slug             string         `json:"slug"`
	Subtitle         string         `json:"subtitle"`
	Description      string         `json:"description"`
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
			MemberID:  row.MemberID.String(),
			Name:      row.MemberName.String,
			Nim:       row.MemberNim.String,
			PhotoUrl:  row.PhotoUrl.String,
			RoleTitle: row.RoleTitle.String,
		}
		if orgTreeCoreSlugs[strings.ToLower(row.DivisionSlug.String)] {
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
		if orgTreeCoreSlugs[strings.ToLower(row.Slug.String)] {
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
			Responsibilities: responsibilities,
		}
		if row.CoordinatorID.Valid {
			division.Coordinator = &OrgTreePerson{
				MemberID:  row.CoordinatorID.UUID.String(),
				Name:      row.CoordinatorName.String,
				Nim:       row.CoordinatorNim.String,
				PhotoUrl:  row.CoordinatorPhotoUrl.String,
				RoleTitle: "Koordinator",
			}
		} else if person, ok := fallbackCoordinator[row.ID.String()]; ok {
			division.Coordinator = &person
		}

		resp.Divisions = append(resp.Divisions, division)
	}

	return resp, nil
}
