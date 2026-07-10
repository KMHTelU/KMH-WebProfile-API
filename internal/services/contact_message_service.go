package services

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *Service) CreateContactMessageService(req requests.CreateContactMessageRequest, c fiber.Ctx) *fiber.Error {
	params := generated.InsertContactMessageParams{
		ID:      uuid.New(),
		Name:    sql.NullString{String: req.Name, Valid: req.Name != ""},
		Email:   sql.NullString{String: req.Email, Valid: req.Email != ""},
		Subject: sql.NullString{String: req.Subject, Valid: req.Subject != ""},
		Message: sql.NullString{String: req.Message, Valid: req.Message != ""},
		IsRead:  sql.NullBool{Bool: false, Valid: true},
	}

	if _, err := s.Repository.CreateContactMessage(params, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create contact message")
	}

	return nil
}

func (s *Service) GetContactMessageByIDService(id uuid.UUID, c fiber.Ctx) (generated.ContactMessage, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return generated.ContactMessage{}, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	message, err1 := s.Repository.GetContactMessageByID(id, c)
	if err1 != nil {
		return generated.ContactMessage{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to get contact message")
	}
	return message, nil
}

func (s *Service) GetPaginatedContactMessagesService(limit, offset int32, c fiber.Ctx) ([]generated.ContactMessage, *fiber.Error) {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	messages, err1 := s.Repository.ListContactMessages(generated.ListContactMessagesParams{
		Limit:  limit,
		Offset: offset,
	}, c)
	if err1 != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get contact messages")
	}
	return messages, nil
}

func (s *Service) MarkContactMessageAsReadService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if _, err := s.Repository.MarkContactMessageAsRead(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to mark contact message as read")
	}
	return nil
}

func (s *Service) DeleteContactMessageService(id uuid.UUID, c fiber.Ctx) *fiber.Error {
	claim, err := s.TokenCleaner.GetCleanToken(c)
	if err != nil || claim == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := s.Repository.DeleteContactMessage(id, c); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete contact message")
	}
	return nil
}
