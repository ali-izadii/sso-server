package handler

import (
	"net/http"

	"github.com/ali/sso-server/internal/model"
	"github.com/ali/sso-server/pkg/logger"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	// TODO: add auth service
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "User registration data"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req model.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("failed to bind register request", "error", err)
		return badRequest(c, "invalid request body")
	}

	// TODO: validate request
	// TODO: check if user exists
	// TODO: hash password
	// TODO: create user

	logger.Info("user registered", "email", req.Email)

	return c.JSON(http.StatusCreated, model.UserResponse{
		Email: req.Email,
		Name:  req.Name,
	})
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("failed to bind login request", "error", err)
		return badRequest(c, "invalid request body")
	}

	// TODO: validate request
	// TODO: find user by email
	// TODO: verify password
	// TODO: generate tokens
	// TODO: create session

	logger.Info("user logged in", "email", req.Email)

	return c.JSON(http.StatusOK, model.TokenResponse{
		AccessToken:  "access_token_placeholder",
		RefreshToken: "refresh_token_placeholder",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} model.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req model.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("failed to bind refresh request", "error", err)
		return badRequest(c, "invalid request body")
	}

	// TODO: validate refresh token
	// TODO: find session
	// TODO: generate new tokens
	// TODO: rotate refresh token

	logger.Info("token refreshed")

	return c.JSON(http.StatusOK, model.TokenResponse{
		AccessToken:  "new_access_token_placeholder",
		RefreshToken: "new_refresh_token_placeholder",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

// Logout godoc
// @Summary Logout user
// @Tags auth
// @Security BearerAuth
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	// TODO: get user from context
	// TODO: invalidate session

	logger.Info("user logged out")

	return c.NoContent(http.StatusNoContent)
}
