package services

import (
	"time"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/mailer"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/repositories"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/cloudinary/cloudinary-go/v2"
)

// PasswordResetConfig mengatur perilaku alur lupa password.
type PasswordResetConfig struct {
	// FrontendURL adalah asal halaman reset password yang dibuka pengguna.
	FrontendURL string
	// TokenTTL menentukan berapa lama tautan reset berlaku.
	TokenTTL time.Duration
	// MaxPerHour membatasi jumlah permintaan reset per akun per jam.
	MaxPerHour int
}

type Service struct {
	// Add necessary fields here, e.g., repository references
	Repository    *repositories.Repository
	TokenCleaner  *utils.TokenCleaner
	Cloudinary    *cloudinary.Cloudinary
	Mailer        *mailer.Mailer
	PasswordReset PasswordResetConfig
}

func InitializeService(
	repo *repositories.Repository,
	tokenCleaner *utils.TokenCleaner,
	cloudinary *cloudinary.Cloudinary,
	mail *mailer.Mailer,
	passwordReset PasswordResetConfig,
) *Service {
	return &Service{
		Repository:    repo,
		TokenCleaner:  tokenCleaner,
		Cloudinary:    cloudinary,
		Mailer:        mail,
		PasswordReset: passwordReset,
	}
}
