package repositories

// Pencarian baris berdasarkan kunci yang dikenal manusia. Dipakai alur import
// untuk mengubah isi kolom Excel seperti NIM atau slug menjadi UUID, dan untuk
// mendeteksi data yang sudah ada sebelum menyisipkan baris baru.

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

func (r *Repository) GetMemberByNIM(nim string, c fiber.Ctx) (generated.Member, error) {
	return r.Queries.GetMemberByNIM(c, utils.NullString(nim))
}

func (r *Repository) GetDivisionBySlug(slug string, c fiber.Ctx) (generated.Division, error) {
	return r.Queries.GetDivisionBySlug(c, utils.NullString(slug))
}

func (r *Repository) GetEventBySlug(slug string, c fiber.Ctx) (generated.Event, error) {
	return r.Queries.GetEventBySlug(c, utils.NullString(slug))
}

func (r *Repository) GetRoleByName(name string, c fiber.Ctx) (generated.Role, error) {
	return r.Queries.GetRoleByName(c, utils.NullString(name))
}
