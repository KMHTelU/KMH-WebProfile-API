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
	return params.Limit, params.Offset
}
