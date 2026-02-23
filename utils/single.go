package utils

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type SingleParams struct {
	ID uuid.UUID `uri:"id"`
}

func GetSingleParams(c fiber.Ctx) uuid.UUID {
	var params SingleParams
	if err := c.Bind().URI(&params); err != nil {
		return uuid.Nil
	}
	return params.ID
}
