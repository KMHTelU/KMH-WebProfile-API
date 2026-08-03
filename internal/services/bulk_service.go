package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) BulkCreateMembersService(req requests.BulkCreateMembersRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.CreateMemberRequest) (uuid.UUID, *fiber.Error) {
		return s.createMember(item, c)
	})
	s.logBulk(claim, "Bulk Create Members", "Member", report, c)
	return report, nil
}

func (s *Service) BulkUpdateMembersService(req requests.BulkUpdateMembersRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.BulkUpdateMemberItem) (uuid.UUID, *fiber.Error) {
		if err := s.updateMember(item.ID, item.UpdateMemberRequest, c); err != nil {
			return uuid.Nil, err
		}
		return item.ID, nil
	})
	s.logBulk(claim, "Bulk Update Members", "Member", report, c)
	return report, nil
}

func (s *Service) BulkCreateDivisionsService(req requests.BulkCreateDivisionsRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.CreateDivisionRequest) (uuid.UUID, *fiber.Error) {
		return s.createDivision(item, c)
	})
	s.logBulk(claim, "Bulk Create Divisions", "Division", report, c)
	return report, nil
}

func (s *Service) BulkUpdateDivisionsService(req requests.BulkUpdateDivisionsRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.BulkUpdateDivisionItem) (uuid.UUID, *fiber.Error) {
		if err := s.updateDivision(item.ID, item.UpdateDivisionRequest, c); err != nil {
			return uuid.Nil, err
		}
		return item.ID, nil
	})
	s.logBulk(claim, "Bulk Update Divisions", "Division", report, c)
	return report, nil
}

func (s *Service) BulkCreateMemberDivisionsService(req requests.BulkCreateMemberDivisionsRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.CreateMemberDivisionRequest) (uuid.UUID, *fiber.Error) {
		return s.createMemberDivision(item, c)
	})
	s.logBulk(claim, "Bulk Create Member Divisions", "Member Division", report, c)
	return report, nil
}

func (s *Service) BulkUpdateMemberDivisionsService(req requests.BulkUpdateMemberDivisionsRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.BulkUpdateMemberDivisionItem) (uuid.UUID, *fiber.Error) {
		if err := s.updateMemberDivision(item.ID, item.UpdateMemberDivisionRequest, c); err != nil {
			return uuid.Nil, err
		}
		return item.ID, nil
	})
	s.logBulk(claim, "Bulk Update Member Divisions", "Member Division", report, c)
	return report, nil
}

func (s *Service) BulkCreateEventsService(req requests.BulkCreateEventsRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.CreateEventRequest) (uuid.UUID, *fiber.Error) {
		return s.createEvent(item, claim.UserID, c)
	})
	s.logBulk(claim, "Bulk Create Events", "Event", report, c)
	return report, nil
}

func (s *Service) BulkUpdateEventsService(req requests.BulkUpdateEventsRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.BulkUpdateEventItem) (uuid.UUID, *fiber.Error) {
		if err := s.updateEvent(item.ID, item.UpdateEventRequest, c); err != nil {
			return uuid.Nil, err
		}
		return item.ID, nil
	})
	s.logBulk(claim, "Bulk Update Events", "Event", report, c)
	return report, nil
}

func (s *Service) BulkCreateGalleriesService(req requests.BulkCreateGalleriesRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.CreateGalleryRequest) (uuid.UUID, *fiber.Error) {
		gallery, err := s.createGallery(item, c)
		if err != nil {
			return uuid.Nil, err
		}
		return gallery.ID, nil
	})
	s.logBulk(claim, "Bulk Create Galleries", "Gallery", report, c)
	return report, nil
}

func (s *Service) BulkUpdateGalleriesService(req requests.BulkUpdateGalleriesRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.BulkUpdateGalleryItem) (uuid.UUID, *fiber.Error) {
		if err := s.updateGallery(item.ID, item.UpdateGalleryRequest, c); err != nil {
			return uuid.Nil, err
		}
		return item.ID, nil
	})
	s.logBulk(claim, "Bulk Update Galleries", "Gallery", report, c)
	return report, nil
}

func (s *Service) BulkCreateGalleryItemsService(req requests.BulkCreateGalleryItemsRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.CreateGalleryItemRequest) (uuid.UUID, *fiber.Error) {
		return s.createGalleryItem(item, c)
	})
	s.logBulk(claim, "Bulk Create Gallery Items", "Gallery Item", report, c)
	return report, nil
}

func (s *Service) BulkUpdateGalleryItemsService(req requests.BulkUpdateGalleryItemsRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.BulkUpdateGalleryItemEntry) (uuid.UUID, *fiber.Error) {
		if err := s.updateGalleryItem(item.ID, item.UpdateGalleryItemRequest, c); err != nil {
			return uuid.Nil, err
		}
		return item.ID, nil
	})
	s.logBulk(claim, "Bulk Update Gallery Items", "Gallery Item", report, c)
	return report, nil
}

func (s *Service) BulkCreateHomepageBannersService(req requests.BulkCreateHomepageBannersRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.HomepageBannerJSONRequest) (uuid.UUID, *fiber.Error) {
		return s.createHomepageBanner(item, c)
	})
	s.logBulk(claim, "Bulk Create Homepage Banners", "Homepage Banner", report, c)
	return report, nil
}

func (s *Service) BulkUpdateHomepageBannersService(req requests.BulkUpdateHomepageBannersRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.BulkUpdateHomepageBannerItem) (uuid.UUID, *fiber.Error) {
		if err := s.updateHomepageBanner(item.ID, item.HomepageBannerJSONRequest, c); err != nil {
			return uuid.Nil, err
		}
		return item.ID, nil
	})
	s.logBulk(claim, "Bulk Update Homepage Banners", "Homepage Banner", report, c)
	return report, nil
}

func (s *Service) BulkCreateRolesService(req requests.BulkCreateRolesRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.CreateRoleRequest) (uuid.UUID, *fiber.Error) {
		return s.createRole(item, c)
	})
	s.logBulk(claim, "Bulk Create Roles", "Role", report, c)
	return report, nil
}

func (s *Service) BulkUpdateRolesService(req requests.BulkUpdateRolesRequest, c fiber.Ctx) (*utils.BulkReport, *fiber.Error) {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return nil, ferr
	}
	report := runBulk(req.Items, func(_ int, item requests.BulkUpdateRoleItem) (uuid.UUID, *fiber.Error) {
		if err := s.updateRole(item.ID, item.UpdateRoleRequest, c); err != nil {
			return uuid.Nil, err
		}
		return item.ID, nil
	})
	s.logBulk(claim, "Bulk Update Roles", "Role", report, c)
	return report, nil
}
