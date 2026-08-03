package services

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/mailer"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
)

// ForgotPasswordService membuat token reset dan mengirimkannya lewat email.
//
// Fungsi ini selalu mengembalikan nil selama tidak ada kegagalan sistem, baik
// email terdaftar maupun tidak. Membedakan keduanya akan membuat endpoint ini
// bisa dipakai memetakan alamat email mana yang punya akun.
func (s *Service) ForgotPasswordService(req requests.ForgotPasswordRequest, c fiber.Ctx) *fiber.Error {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.Repository.GetUserByEmail(email, c)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("Gagal mencari user saat forgot password: %v", err)
		}
		return nil
	}

	if user.IsActive.Valid && !user.IsActive.Bool {
		log.Infof("Permintaan reset password untuk akun nonaktif: %s", user.ID)
		return nil
	}

	recent, err := s.Repository.CountRecentPasswordResetTokens(user.ID, time.Now().Add(-time.Hour), c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process password reset request")
	}
	if recent >= int64(s.PasswordReset.MaxPerHour) {
		return fiber.NewError(fiber.StatusTooManyRequests, "Too many password reset requests, please try again later")
	}

	// Token lama dibatalkan supaya hanya tautan terbaru yang bisa dipakai.
	if err := s.Repository.InvalidateUserPasswordResetTokens(user.ID, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process password reset request")
	}

	token, err := utils.GenerateResetToken()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process password reset request")
	}

	if err := s.Repository.CreatePasswordResetToken(generated.InsertPasswordResetTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: utils.HashResetToken(token),
		ExpiresAt: time.Now().Add(s.PasswordReset.TokenTTL),
	}, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process password reset request")
	}

	message, err := s.Mailer.BuildResetPassword(user.Email.String, mailer.ResetPasswordData{
		UserName:         user.Name.String,
		ResetURL:         s.buildResetURL(token),
		ExpiresInMinutes: int(s.PasswordReset.TokenTTL.Minutes()),
	})
	if err != nil {
		log.Errorf("Gagal menyusun email reset password: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process password reset request")
	}
	s.Mailer.SendAsync(message)

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(user.ID),
		Action:    utils.NullString("Request Password Reset"),
		Entity:    utils.NullString("User"),
		EntityID:  utils.NullUUID(user.ID),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

// ResetPasswordService menukar token yang sah dengan password baru.
func (s *Service) ResetPasswordService(req requests.ResetPasswordRequest, c fiber.Ctx) *fiber.Error {
	invalidToken := fiber.NewError(fiber.StatusBadRequest, "Reset token is invalid or has expired")

	record, err := s.Repository.GetPasswordResetTokenByHash(utils.HashResetToken(req.Token), c)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("Gagal membaca token reset password: %v", err)
		}
		return invalidToken
	}
	if record.UsedAt.Valid {
		return invalidToken
	}
	if time.Now().After(record.ExpiresAt) {
		return invalidToken
	}

	user, err := s.Repository.GetUserByID(record.UserID, c)
	if err != nil {
		return invalidToken
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to reset password")
	}
	if err := s.Repository.UpdateUserPassword(record.UserID, hashed, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to reset password")
	}
	if err := s.Repository.MarkPasswordResetTokenUsed(record.ID, c); err != nil {
		log.Errorf("Password berhasil diubah tetapi token gagal ditandai terpakai: %v", err)
	}

	s.notifyPasswordChanged(user.Email.String, user.Name.String)

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(record.UserID),
		Action:    utils.NullString("Reset Password"),
		Entity:    utils.NullString("User"),
		EntityID:  utils.NullUUID(record.UserID),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

// ChangePasswordService mengubah password pengguna yang sedang login.
func (s *Service) ChangePasswordService(req requests.ChangePasswordRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	user, err := s.Repository.GetUserByID(claim.UserID, c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if !utils.CheckPassword(user.PasswordHash.String, req.OldPassword) {
		return fiber.NewError(fiber.StatusBadRequest, "Current password is incorrect")
	}
	if req.OldPassword == req.NewPassword {
		return fiber.NewError(fiber.StatusBadRequest, "New password must be different from the current password")
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to change password")
	}
	if err := s.Repository.UpdateUserPassword(user.ID, hashed, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to change password")
	}

	// Tautan reset yang masih menganggur ikut dibatalkan: setelah pemilik akun
	// mengubah password sendiri, tautan lama tidak boleh lagi bisa dipakai.
	if err := s.Repository.InvalidateUserPasswordResetTokens(user.ID, c); err != nil {
		log.Errorf("Gagal membatalkan token reset setelah ganti password: %v", err)
	}

	s.notifyPasswordChanged(user.Email.String, user.Name.String)

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(user.ID),
		Action:    utils.NullString("Change Password"),
		Entity:    utils.NullString("User"),
		EntityID:  utils.NullUUID(user.ID),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

func (s *Service) buildResetURL(token string) string {
	base := strings.TrimSuffix(s.PasswordReset.FrontendURL, "/")
	return base + "/reset-password?token=" + url.QueryEscape(token)
}

func (s *Service) notifyPasswordChanged(email, name string) {
	if email == "" {
		return
	}
	message, err := s.Mailer.BuildPasswordChanged(email, mailer.PasswordChangedData{
		UserName:  name,
		ChangedAt: time.Now().Format("2 January 2006 15:04 MST"),
	})
	if err != nil {
		log.Errorf("Gagal menyusun email notifikasi perubahan password: %v", err)
		return
	}
	s.Mailer.SendAsync(message)
}
