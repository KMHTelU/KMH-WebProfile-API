package utils

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// Status yang mungkin muncul pada tiap item dalam laporan bulk.
const (
	BulkStatusSuccess = "success"
	BulkStatusFailed  = "failed"
)

// MaxBulkItems membatasi jumlah item yang boleh diproses dalam satu permintaan
// bulk maupun satu berkas import.
const MaxBulkItems = 500

// BulkItemResult adalah hasil pemrosesan satu item.
type BulkItemResult struct {
	Index  int    `json:"index"`
	Status string `json:"status"`
	ID     string `json:"id,omitempty"`
	// Row diisi pada alur import sebagai nomor baris berkas yang sebenarnya,
	// sehingga admin bisa langsung membuka baris tersebut di Excel.
	Row   int    `json:"row,omitempty"`
	Error string `json:"error,omitempty"`
}

// BulkReport merangkum hasil seluruh item pada satu operasi bulk atau import.
type BulkReport struct {
	Total     int              `json:"total"`
	Succeeded int              `json:"succeeded"`
	Failed    int              `json:"failed"`
	Results   []BulkItemResult `json:"results"`
}

func NewBulkReport(total int) *BulkReport {
	return &BulkReport{
		Total:   total,
		Results: make([]BulkItemResult, 0, total),
	}
}

func (r *BulkReport) AddSuccess(index int, id uuid.UUID) {
	r.Succeeded++
	r.Results = append(r.Results, BulkItemResult{
		Index:  index,
		Status: BulkStatusSuccess,
		ID:     id.String(),
	})
}

func (r *BulkReport) AddFailure(index int, message string) {
	r.Failed++
	r.Results = append(r.Results, BulkItemResult{
		Index:  index,
		Status: BulkStatusFailed,
		Error:  message,
	})
}

// AddRowSuccess dan AddRowFailure dipakai alur import, yang perlu menyebut nomor
// baris berkas selain indeks item.
func (r *BulkReport) AddRowSuccess(index, row int, id uuid.UUID) {
	r.Succeeded++
	r.Results = append(r.Results, BulkItemResult{
		Index:  index,
		Row:    row,
		Status: BulkStatusSuccess,
		ID:     id.String(),
	})
}

func (r *BulkReport) AddRowFailure(index, row int, message string) {
	r.Failed++
	r.Results = append(r.Results, BulkItemResult{
		Index:  index,
		Row:    row,
		Status: BulkStatusFailed,
		Error:  message,
	})
}

// RespondWithBulkReport memilih kode status berdasarkan hasil keseluruhan:
// seluruhnya berhasil memakai successStatus, sebagian berhasil memakai
// 207 Multi-Status, dan seluruhnya gagal memakai 400.
func RespondWithBulkReport(c fiber.Ctx, message string, report *BulkReport, successStatus int) error {
	switch {
	case report.Failed == 0:
		return c.Status(successStatus).JSON(APIResponse{
			Status:  "success",
			Message: message,
			Data:    report,
		})
	case report.Succeeded == 0:
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Status:  "error",
			Message: message + " (semua item gagal)",
			Data:    report,
		})
	default:
		return c.Status(fiber.StatusMultiStatus).JSON(APIResponse{
			Status:  "partial",
			Message: message + " (sebagian item gagal)",
			Data:    report,
		})
	}
}
