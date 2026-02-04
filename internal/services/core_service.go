package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/repositories"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
)

type Service struct {
	// Add necessary fields here, e.g., repository references
	Repository   *repositories.Repository
	TokenCleaner *utils.TokenCleaner
}

func InitializeService(repo *repositories.Repository, tokenCleaner *utils.TokenCleaner) *Service {
	return &Service{
		Repository:   repo,
		TokenCleaner: tokenCleaner,
	}
}
