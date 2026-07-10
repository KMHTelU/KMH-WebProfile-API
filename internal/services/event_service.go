package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateEventService(req requests.CreateEventRequest, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.InsertEventParams{
		ID:              uuid.New(),
		Title:           sql.NullString{String: req.Title, Valid: req.Title != ""},
		Slug:            sql.NullString{String: req.Slug, Valid: req.Slug != ""},
		Description:     sql.NullString{String: req.Description, Valid: req.Description != ""},
		EventType:       sql.NullString{String: req.EventType, Valid: req.EventType != ""},
		StartTime:       sql.NullTime{Time: req.StartTime, Valid: !req.StartTime.IsZero()},
		EndTime:         sql.NullTime{Time: req.EndTime, Valid: !req.EndTime.IsZero()},
		Location:        sql.NullString{String: req.Location, Valid: req.Location != ""},
		GoogleMapsUrl:   sql.NullString{String: req.GoogleMapsUrl, Valid: req.GoogleMapsUrl != ""},
		RegistrationUrl: sql.NullString{String: req.RegistrationUrl, Valid: req.RegistrationUrl != ""},
		CoverMediaID:    uuid.NullUUID{UUID: req.CoverMediaID, Valid: req.CoverMediaID != uuid.Nil},
		Status:          sql.NullString{String: req.Status, Valid: req.Status != ""},
		IsPublished:     sql.NullBool{Bool: req.IsPublished, Valid: true},
		CreatedBy:       uuid.NullUUID{UUID: claim.UserID, Valid: true},
	}

	if err := s.Repository.CreateEvent(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create event")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Create Event", Valid: true},
		Entity:    sql.NullString{String: "Event", Valid: true},
		EntityID:  uuid.NullUUID{UUID: params.ID, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
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
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	params := generated.UpdateEventParams{
		ID:              id,
		Title:           sql.NullString{String: req.Title, Valid: req.Title != ""},
		Slug:            sql.NullString{String: req.Slug, Valid: req.Slug != ""},
		Description:     sql.NullString{String: req.Description, Valid: req.Description != ""},
		EventType:       sql.NullString{String: req.EventType, Valid: req.EventType != ""},
		StartTime:       sql.NullTime{Time: req.StartTime, Valid: !req.StartTime.IsZero()},
		EndTime:         sql.NullTime{Time: req.EndTime, Valid: !req.EndTime.IsZero()},
		Location:        sql.NullString{String: req.Location, Valid: req.Location != ""},
		GoogleMapsUrl:   sql.NullString{String: req.GoogleMapsUrl, Valid: req.GoogleMapsUrl != ""},
		RegistrationUrl: sql.NullString{String: req.RegistrationUrl, Valid: req.RegistrationUrl != ""},
		CoverMediaID:    uuid.NullUUID{UUID: req.CoverMediaID, Valid: req.CoverMediaID != uuid.Nil},
		Status:          sql.NullString{String: req.Status, Valid: req.Status != ""},
		IsPublished:     sql.NullBool{Bool: req.IsPublished, Valid: true},
	}

	if err := s.Repository.UpdateEvent(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update event")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Update Event", Valid: true},
		Entity:    sql.NullString{String: "Event", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}

func (s *Service) DeleteEventService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteEvent(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete event")
	}

	s.Repository.InsertLog(generated.InsertActivityLogParams{
		ID:        uuid.New(),
		UserID:    uuid.NullUUID{UUID: claim.UserID, Valid: true},
		Action:    sql.NullString{String: "Delete Event", Valid: true},
		Entity:    sql.NullString{String: "Event", Valid: true},
		EntityID:  uuid.NullUUID{UUID: id, Valid: true},
		IpAddress: sql.NullString{String: c.IP(), Valid: true},
		UserAgent: sql.NullString{String: c.UserAgent(), Valid: true},
	}, c)

	return nil
}
