package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/repositories"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/cloudinary/cloudinary-go/v2"
)

type Service struct {
	// Add necessary fields here, e.g., repository references
	Repository   *repositories.Repository
	TokenCleaner *utils.TokenCleaner
	Cloudinary   *cloudinary.Cloudinary
}

func InitializeService(repo *repositories.Repository, tokenCleaner *utils.TokenCleaner, cloudinary *cloudinary.Cloudinary) *Service {
	return &Service{
		Repository:   repo,
		TokenCleaner: tokenCleaner,
		Cloudinary:   cloudinary,
	}
}
