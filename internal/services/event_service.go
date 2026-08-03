package services

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateEventService(req requests.CreateEventRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	id, ferr := s.createEvent(req, claim.UserID, c)
	if ferr != nil {
		return ferr
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Create Event"),
		Entity:    utils.NullString("Event"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

func (s *Service) GetEventByIDService(id uuid.UUID, c fiber.Ctx) (generated.GetEventByIDRow, *fiber.Error) {
	event, err := s.Repository.GetEventByID(id, c)
	if err != nil {
		return generated.GetEventByIDRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get event")
	}
	return event, nil
}

func (s *Service) GetPaginatedEventsService(limit, offset int32, c fiber.Ctx) ([]generated.ListEventsRow, *fiber.Error) {
	events, err := s.Repository.ListEvents(generated.ListEventsParams{
		Limit:  limit,
		Offset: offset,
	}, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get events")
	}
	return events, nil
}

func (s *Service) UpdateEventService(id uuid.UUID, req requests.UpdateEventRequest, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	if ferr := s.updateEvent(id, req, c); ferr != nil {
		return ferr
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Update Event"),
		Entity:    utils.NullString("Event"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}

func (s *Service) DeleteEventService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, ferr := s.requireAuth(c)
	if ferr != nil {
		return ferr
	}

	if err := s.Repository.DeleteEvent(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete event")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    utils.NullUUID(claim.UserID),
		Action:    utils.NullString("Delete Event"),
		Entity:    utils.NullString("Event"),
		EntityID:  utils.NullUUID(id),
		IpAddress: utils.NullString(c.IP()),
		UserAgent: utils.NullString(c.UserAgent()),
	}, c)

	return nil
}
