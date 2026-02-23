package handlers

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/entities"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/requests"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (h *Handler) CreateUser(c fiber.Ctx) error {
	// Implementation for creating a user
	var request requests.CreateUserRequest
	if err := c.Bind().JSON(request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}

	if err := h.Service.CreateUserService(request, c); err != nil {
		//if failed, then must check repository.InsertUser and repository.InsertLog
		return utils.RespondWithError(c, err.Code, err.Message)
	}

	return utils.RespondWithCreated(c, "User created successfully", nil)
}

func (h *Handler) GetAllUsers(c fiber.Ctx) error {
	limit, offset := utils.GetPaginationParams(c)
	users, err := h.Service.GetAllUsersService(limit, offset, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	var datas []entities.User
	for _, user := range users {
		var role entities.Role
		role.ID = user.RoleID
		role.Name = user.Name_2.String
		datas = append(datas, entities.User{
			ID:       user.ID,
			Name:     user.Name.String,
			Email:    user.Email.String,
			RoleID:   user.RoleID,
			UserRole: role,
		})
	}
	return utils.RespondWithOK(c, "Users retrieved successfully", datas)
}

func (h *Handler) GetUserByID(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	user, err := h.Service.GetUserByIDService(id, c)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "User retrieved successfully", user)
}

func (h *Handler) UpdateUser(c fiber.Ctx) error {
	// Implementation for updating a user
	var request requests.UpdateUserRequest
	if err := c.Bind().JSON(request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if err := h.Service.UpdateUserService(id, request, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "User updated successfully", nil)
}

func (h *Handler) DeleteUser(c fiber.Ctx) error {
	id := utils.GetSingleParams(c)
	if id == uuid.Nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if err := h.Service.DeleteUserService(id, c); err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "User deleted successfully", nil)
}
