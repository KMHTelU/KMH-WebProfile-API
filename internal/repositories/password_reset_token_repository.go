package repositories

import (
	"time"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (r *Repository) CreatePasswordResetToken(params generated.InsertPasswordResetTokenParams, c fiber.Ctx) error {
	_, err := r.Queries.InsertPasswordResetToken(c, params)
	return err
}

func (r *Repository) GetPasswordResetTokenByHash(tokenHash string, c fiber.Ctx) (generated.PasswordResetToken, error) {
	return r.Queries.GetPasswordResetTokenByHash(c, tokenHash)
}

func (r *Repository) MarkPasswordResetTokenUsed(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.MarkPasswordResetTokenUsed(c, id)
}

func (r *Repository) InvalidateUserPasswordResetTokens(userID uuid.UUID, c fiber.Ctx) error {
	return r.Queries.InvalidateUserPasswordResetTokens(c, userID)
}

func (r *Repository) CountRecentPasswordResetTokens(userID uuid.UUID, since time.Time, c fiber.Ctx) (int64, error) {
	return r.Queries.CountRecentPasswordResetTokens(c, generated.CountRecentPasswordResetTokensParams{
		UserID:    userID,
		CreatedAt: utils.NullTime(since),
	})
}

func (r *Repository) UpdateUserPassword(userID uuid.UUID, passwordHash string, c fiber.Ctx) error {
	return r.Queries.UpdateUserPassword(c, generated.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: utils.NullString(passwordHash),
	})
}
