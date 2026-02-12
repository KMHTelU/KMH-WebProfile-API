package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/entities"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

func (s *Service) AuthenticateUserService(request requests.AuthenticateUserRequest, c fiber.Ctx) (entities.Auth, *fiber.Error) {
	// Implementation for authenticating a user
	user, err := s.Repository.GetUserByEmail(request.Email, c)
	if err != nil {
		return entities.Auth{}, fiber.NewError(fiber.StatusBadRequest, "Invalid email or password")
	}

	if !utils.CheckPassword(user.PasswordHash.String, request.Password) {
		return entities.Auth{}, fiber.NewError(fiber.StatusBadRequest, "Invalid email or password")
	}

	accessToken, refreshToken, expiresAt, err := utils.GenerateJWT(user.ID, s.TokenCleaner.AccessSecretKey, s.TokenCleaner.RefreshSecretKey)
	if err != nil {
		return entities.Auth{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userResponse := entities.User{
		ID:     user.ID,
		Name:   user.Name.String,
		Email:  user.Email.String,
		RoleID: user.RoleID,
		UserRole: entities.Role{
			ID:          user.ID_2,
			Name:        user.Name_2.String,
			Description: user.Description.String,
		},
		IsActive:    user.IsActive.Bool,
		LastLoginAt: user.LastLoginAt.Time,
	}

	auth := entities.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         userResponse,
	}

	return auth, nil
}

func (s *Service) RefreshTokenService(refreshToken string, c fiber.Ctx) (entities.Auth, *fiber.Error) {
	// Implementation for refreshing JWT token
	claims, err := utils.ValidateJWT(refreshToken, s.TokenCleaner.RefreshSecretKey)
	if err != nil {
		return entities.Auth{}, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	user, err := s.Repository.GetUserByID(claims.UserID, c)
	if err != nil {
		return entities.Auth{}, fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}
	accessToken, newRefreshToken, expiresAt, err := utils.GenerateJWT(user.ID, s.TokenCleaner.AccessSecretKey, s.TokenCleaner.RefreshSecretKey)
	if err != nil {
		return entities.Auth{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	auth := entities.Auth{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}
	return auth, nil
}
