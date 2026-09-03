package utils

import "github.com/gofiber/fiber/v3"

type PaginationParams struct {
	Limit  int32 `query:"limit"`
	Offset int32 `query:"start"`
}

func GetPaginationParams(c fiber.Ctx) (int32, int32) {
	var params PaginationParams
	if err := c.Bind().Query(&params); err != nil {
		return 10, 0 // Default values if parsing fails
	}
	// Terapkan default & batas aman walau parsing sukses.
	// Tanpa ini, "limit" yang tidak dikirim = 0 => LIMIT 0 => hasil kosong.
	if params.Limit <= 0 {
		params.Limit = 10
	}
	// Batas atas dinaikkan supaya dashboard bisa memuat seluruh data (anggota
	// sudah 140+); tetap dibatasi agar query tak terkendali tidak mungkin.
	if params.Limit > 1000 {
		params.Limit = 1000
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	return params.Limit, params.Offset
}
