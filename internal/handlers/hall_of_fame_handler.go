package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// bindJSON menyatukan pola bind + validasi yang dipakai handler HOF.
func bindJSON[T any](c fiber.Ctx, out *T) error {
	if err := c.Bind().JSON(out); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	return nil
}

func requireID(c fiber.Ctx) (uuid.UUID, bool) {
	id := utils.GetSingleParams(c)
	return id, id != uuid.Nil
}

// GetHallOfFameHandler publik: seluruh arsip untuk museum 3D.
func (h *Handler) GetHallOfFameHandler(c fiber.Ctx) error {
	data, err := h.Service.GetHallOfFameService(c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Hall of fame retrieved successfully", data)
}

// ── Generations ──

func (h *Handler) CreateHofGenerationHandler(c fiber.Ctx) error {
	var req requests.HofGenerationRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}
	row, ferr := h.Service.CreateHofGenerationService(req, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithCreated(c, "HOF generation created successfully", row)
}

func (h *Handler) UpdateHofGenerationHandler(c fiber.Ctx) error {
	id, ok := requireID(c)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.HofGenerationRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}
	if ferr := h.Service.UpdateHofGenerationService(id, req, c); ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "HOF generation updated successfully", nil)
}

func (h *Handler) DeleteHofGenerationHandler(c fiber.Ctx) error {
	id, ok := requireID(c)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if ferr := h.Service.DeleteHofGenerationService(id, c); ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "HOF generation deleted successfully", nil)
}

// ── People ──

func (h *Handler) CreateHofPersonHandler(c fiber.Ctx) error {
	var req requests.HofPersonRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}
	row, ferr := h.Service.CreateHofPersonService(req, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithCreated(c, "HOF person created successfully", row)
}

func (h *Handler) UpdateHofPersonHandler(c fiber.Ctx) error {
	id, ok := requireID(c)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.HofPersonRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}
	if ferr := h.Service.UpdateHofPersonService(id, req, c); ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "HOF person updated successfully", nil)
}

func (h *Handler) DeleteHofPersonHandler(c fiber.Ctx) error {
	id, ok := requireID(c)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if ferr := h.Service.DeleteHofPersonService(id, c); ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "HOF person deleted successfully", nil)
}

// ── Achievements ──

func (h *Handler) CreateHofAchievementHandler(c fiber.Ctx) error {
	var req requests.HofAchievementRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}
	row, ferr := h.Service.CreateHofAchievementService(req, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithCreated(c, "HOF achievement created successfully", row)
}

func (h *Handler) UpdateHofAchievementHandler(c fiber.Ctx) error {
	id, ok := requireID(c)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.HofAchievementRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}
	if ferr := h.Service.UpdateHofAchievementService(id, req, c); ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "HOF achievement updated successfully", nil)
}

func (h *Handler) DeleteHofAchievementHandler(c fiber.Ctx) error {
	id, ok := requireID(c)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if ferr := h.Service.DeleteHofAchievementService(id, c); ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "HOF achievement deleted successfully", nil)
}

// ── Timeline ──

func (h *Handler) CreateHofTimelineEventHandler(c fiber.Ctx) error {
	var req requests.HofTimelineEventRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}
	row, ferr := h.Service.CreateHofTimelineEventService(req, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithCreated(c, "HOF timeline event created successfully", row)
}

func (h *Handler) UpdateHofTimelineEventHandler(c fiber.Ctx) error {
	id, ok := requireID(c)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.HofTimelineEventRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}
	if ferr := h.Service.UpdateHofTimelineEventService(id, req, c); ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "HOF timeline event updated successfully", nil)
}

func (h *Handler) DeleteHofTimelineEventHandler(c fiber.Ctx) error {
	id, ok := requireID(c)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if ferr := h.Service.DeleteHofTimelineEventService(id, c); ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithOK(c, "HOF timeline event deleted successfully", nil)
}
