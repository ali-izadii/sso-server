package handler

import (
	"net/http"

	"github.com/ali/sso-server/internal/model"
	"github.com/ali/sso-server/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	// TODO: add user service
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetMe godoc
// @Summary Get current user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetMe(c echo.Context) error {
	// TODO: get user ID from JWT context
	// TODO: fetch user from database

	logger.Debug("fetching current user")

	return c.JSON(http.StatusOK, model.UserResponse{
		ID:    uuid.New(),
		Email: "user@example.com",
		Name:  "John Doe",
	})
}

// UpdateMe godoc
// @Summary Update current user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateUserRequest true "Update data"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/users/me [patch]
func (h *UserHandler) UpdateMe(c echo.Context) error {
	var req model.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("failed to bind update request", "error", err)
		return badRequest(c, "invalid request body")
	}

	// TODO: get user ID from JWT context
	// TODO: update user in database

	logger.Info("user updated")

	name := "John Doe"
	if req.Name != nil {
		name = *req.Name
	}

	return c.JSON(http.StatusOK, model.UserResponse{
		ID:    uuid.New(),
		Email: "user@example.com",
		Name:  name,
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Password change data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/users/me/password [put]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("failed to bind password change request", "error", err)
		return badRequest(c, "invalid request body")
	}

	// TODO: validate old password
	// TODO: hash new password
	// TODO: update password in database
	// TODO: invalidate all sessions except current

	logger.Info("password changed")

	return success(c, "password changed successfully")
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
