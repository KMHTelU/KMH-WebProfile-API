package services

import (
	"fmt"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// requireAuth memusatkan pemeriksaan token yang di service lain ditulis berulang.
func (s *Service) requireAuth(c fiber.Ctx) (*utils.Claims, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	return claim, nil
}

// runBulk memproses setiap item secara berurutan dan mengumpulkan hasilnya.
//
// Satu item yang gagal tidak menghentikan item lain: berkas dan payload hasil
// kerja manual hampir selalu punya beberapa baris bermasalah, dan menolak
// seluruh kiriman karena satu kesalahan ketik membuat alur kerja admin menyakitkan.
func runBulk[T any](items []T, handle func(index int, item T) (uuid.UUID, *fiber.Error)) *utils.BulkReport {
	report := utils.NewBulkReport(len(items))
	for index, item := range items {
		id, err := handle(index, item)
		if err != nil {
			report.AddFailure(index, err.Message)
			continue
		}
		report.AddSuccess(index, id)
	}
	return report
}

// logBulk mencatat satu baris ringkasan per operasi. Mencatat per item akan
// membuat activity_logs membengkak tanpa memberi informasi tambahan.
func (s *Service) logBulk(claim *utils.Claims, action, entity string, report *utils.BulkReport, c fiber.Ctx) {
	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:     uuid.New(),
		UserID: utils.NullUUID(claim.UserID),
		Action: utils.NullString(action),
		Entity: utils.NullString(
			fmt.Sprintf("%s (%d berhasil, %d gagal)", entity, report.Succeeded, report.Failed),
		),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)
}
