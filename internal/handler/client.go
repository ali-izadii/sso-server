package handler

import (
	"net/http"

	"github.com/ali/sso-server/internal/model"
	"github.com/ali/sso-server/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	// TODO: add client service
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}

// Create godoc
// @Summary Register a new OAuth client
// @Tags clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateClientRequest true "Client registration data"
// @Success 201 {object} model.ClientResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/clients [post]
func (h *ClientHandler) Create(c echo.Context) error {
	var req model.CreateClientRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("failed to bind client request", "error", err)
		return badRequest(c, "invalid request body")
	}

	// TODO: validate request
	// TODO: generate client secret
	// TODO: create client in database

	clientID := uuid.New()
	clientSecret := uuid.New().String() // placeholder

	logger.Info("client created", "client_id", clientID)

	return c.JSON(http.StatusCreated, model.ClientResponse{
		ID:           clientID,
		Name:         req.Name,
		Secret:       clientSecret,
		RedirectURIs: req.RedirectURIs,
	})
}

// List godoc
// @Summary List all OAuth clients
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.ClientResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/clients [get]
func (h *ClientHandler) List(c echo.Context) error {
	// TODO: get clients from database

	logger.Debug("listing clients")

	return c.JSON(http.StatusOK, []model.ClientResponse{})
}

// Get godoc
// @Summary Get OAuth client by ID
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} model.ClientResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/clients/{id} [get]
func (h *ClientHandler) Get(c echo.Context) error {
	id := c.Param("id")

	clientID, err := uuid.Parse(id)
	if err != nil {
		return badRequest(c, "invalid client id")
	}

	// TODO: get client from database

	logger.Debug("fetching client", "client_id", clientID)

	return c.JSON(http.StatusOK, model.ClientResponse{
		ID:           clientID,
		Name:         "Example Client",
		RedirectURIs: []string{"https://example.com/callback"},
	})
}

// Delete godoc
// @Summary Delete OAuth client
// @Tags clients
// @Security BearerAuth
// @Param id path string true "Client ID"
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/clients/{id} [delete]
func (h *ClientHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	clientID, err := uuid.Parse(id)
	if err != nil {
		return badRequest(c, "invalid client id")
	}

	// TODO: delete client from database

	logger.Info("client deleted", "client_id", clientID)

	return c.NoContent(http.StatusNoContent)
}
