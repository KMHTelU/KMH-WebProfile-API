package handlers

import "github.com/KMHTelU/KMH-WebProfile-API/internal/services"

type Handler struct {
	// Add necessary fields here, e.g., service references
	Service *services.Service
}

func InitializeHandler(svc *services.Service) *Handler {
	return &Handler{
		Service: svc,
	}
}
