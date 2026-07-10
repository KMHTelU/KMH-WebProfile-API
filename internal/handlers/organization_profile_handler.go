package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateOrganizationProfileHandler(c fiber.Ctx) error {
	var req requests.CreateOrganizationProfileRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.CreateOrganizationProfileService(req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithCreated(c, "Organization Profile created successfully", nil)
}

func (h *Handler) GetOrganizationProfileHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	profile, err := h.Service.GetOrganizationProfileService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Organization Profile retrieved successfully", profile)
}

func (h *Handler) UpdateOrganizationProfileHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	var req requests.UpdateOrganizationProfileRequest
	if err := c.Bind().JSON(&req); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.UpdateOrganizationProfileService(id, req, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Organization Profile updated successfully", nil)
}

func (h *Handler) DeleteOrganizationProfileHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.DeleteOrganizationProfileService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Organization Profile deleted successfully", nil)
}

func (h *Handler) UploadAndUpdateOrganizationProfileLogoHandler(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	err := h.Service.UploadOrganizationProfileLogoService(id, c)
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Organization Profile logo uploaded and updated successfully", nil)
}
