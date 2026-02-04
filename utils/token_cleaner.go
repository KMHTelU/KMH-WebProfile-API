package utils

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

type TokenCleaner struct {
	SecretKey []byte
}

func InitializeTokenCleaner(secret string) *TokenCleaner {
	return &TokenCleaner{
		SecretKey: []byte(secret),
	}
}

func (tc *TokenCleaner) GetCleanToken(c fiber.Ctx) (*Claims, error) {
	token := c.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := ValidateJWT(token, tc.SecretKey)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
