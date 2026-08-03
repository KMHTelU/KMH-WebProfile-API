package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
)

// bindBulk membaca body JSON untuk permintaan bulk.
func bindBulk[T any](c fiber.Ctx, request *T) error {
	if err := c.Bind().JSON(request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	return nil
}

// respondBulk menjalankan service bulk lalu mengubah laporannya menjadi respons.
func respondBulk[T any](
	c fiber.Ctx,
	request *T,
	run func(T, fiber.Ctx) (*utils.BulkReport, *fiber.Error),
	message string,
	successStatus int,
) error {
	if err := bindBulk(c, request); err != nil {
		return err
	}
	report, ferr := run(*request, c)
	if ferr != nil {
		return utils.RespondWithError(c, ferr.Code, ferr.Message)
	}
	return utils.RespondWithBulkReport(c, message, report, successStatus)
}

func (h *Handler) BulkCreateMembersHandler(c fiber.Ctx) error {
	var request requests.BulkCreateMembersRequest
	return respondBulk(c, &request, h.Service.BulkCreateMembersService, "Bulk create members processed", fiber.StatusCreated)
}

func (h *Handler) BulkUpdateMembersHandler(c fiber.Ctx) error {
	var request requests.BulkUpdateMembersRequest
	return respondBulk(c, &request, h.Service.BulkUpdateMembersService, "Bulk update members processed", fiber.StatusOK)
}

func (h *Handler) BulkCreateDivisionsHandler(c fiber.Ctx) error {
	var request requests.BulkCreateDivisionsRequest
	return respondBulk(c, &request, h.Service.BulkCreateDivisionsService, "Bulk create divisions processed", fiber.StatusCreated)
}

func (h *Handler) BulkUpdateDivisionsHandler(c fiber.Ctx) error {
	var request requests.BulkUpdateDivisionsRequest
	return respondBulk(c, &request, h.Service.BulkUpdateDivisionsService, "Bulk update divisions processed", fiber.StatusOK)
}

func (h *Handler) BulkCreateMemberDivisionsHandler(c fiber.Ctx) error {
	var request requests.BulkCreateMemberDivisionsRequest
	return respondBulk(c, &request, h.Service.BulkCreateMemberDivisionsService, "Bulk create member divisions processed", fiber.StatusCreated)
}

func (h *Handler) BulkUpdateMemberDivisionsHandler(c fiber.Ctx) error {
	var request requests.BulkUpdateMemberDivisionsRequest
	return respondBulk(c, &request, h.Service.BulkUpdateMemberDivisionsService, "Bulk update member divisions processed", fiber.StatusOK)
}

func (h *Handler) BulkCreateEventsHandler(c fiber.Ctx) error {
	var request requests.BulkCreateEventsRequest
	return respondBulk(c, &request, h.Service.BulkCreateEventsService, "Bulk create events processed", fiber.StatusCreated)
}

func (h *Handler) BulkUpdateEventsHandler(c fiber.Ctx) error {
	var request requests.BulkUpdateEventsRequest
	return respondBulk(c, &request, h.Service.BulkUpdateEventsService, "Bulk update events processed", fiber.StatusOK)
}

func (h *Handler) BulkCreateGalleriesHandler(c fiber.Ctx) error {
	var request requests.BulkCreateGalleriesRequest
	return respondBulk(c, &request, h.Service.BulkCreateGalleriesService, "Bulk create galleries processed", fiber.StatusCreated)
}

func (h *Handler) BulkUpdateGalleriesHandler(c fiber.Ctx) error {
	var request requests.BulkUpdateGalleriesRequest
	return respondBulk(c, &request, h.Service.BulkUpdateGalleriesService, "Bulk update galleries processed", fiber.StatusOK)
}

func (h *Handler) BulkCreateGalleryItemsHandler(c fiber.Ctx) error {
	var request requests.BulkCreateGalleryItemsRequest
	return respondBulk(c, &request, h.Service.BulkCreateGalleryItemsService, "Bulk create gallery items processed", fiber.StatusCreated)
}

func (h *Handler) BulkUpdateGalleryItemsHandler(c fiber.Ctx) error {
	var request requests.BulkUpdateGalleryItemsRequest
	return respondBulk(c, &request, h.Service.BulkUpdateGalleryItemsService, "Bulk update gallery items processed", fiber.StatusOK)
}

func (h *Handler) BulkCreateHomepageBannersHandler(c fiber.Ctx) error {
	var request requests.BulkCreateHomepageBannersRequest
	return respondBulk(c, &request, h.Service.BulkCreateHomepageBannersService, "Bulk create homepage banners processed", fiber.StatusCreated)
}

func (h *Handler) BulkUpdateHomepageBannersHandler(c fiber.Ctx) error {
	var request requests.BulkUpdateHomepageBannersRequest
	return respondBulk(c, &request, h.Service.BulkUpdateHomepageBannersService, "Bulk update homepage banners processed", fiber.StatusOK)
}

func (h *Handler) BulkCreateRolesHandler(c fiber.Ctx) error {
	var request requests.BulkCreateRolesRequest
	return respondBulk(c, &request, h.Service.BulkCreateRolesService, "Bulk create roles processed", fiber.StatusCreated)
}

func (h *Handler) BulkUpdateRolesHandler(c fiber.Ctx) error {
	var request requests.BulkUpdateRolesRequest
	return respondBulk(c, &request, h.Service.BulkUpdateRolesService, "Bulk update roles processed", fiber.StatusOK)
}
