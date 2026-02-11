package utils

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

type TokenCleaner struct {
	AccessSecretKey  []byte
	RefreshSecretKey []byte
}

func InitializeTokenCleaner(accessSecret string, refreshSecret string) *TokenCleaner {
	return &TokenCleaner{
		AccessSecretKey:  []byte(accessSecret),
		RefreshSecretKey: []byte(refreshSecret),
	}
}

func (tc *TokenCleaner) GetCleanToken(c fiber.Ctx) (*Claims, error) {
	token := c.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := ValidateJWT(token, tc.AccessSecretKey)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
