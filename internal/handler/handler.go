package handler

import "github.com/labstack/echo/v4"

type Handler struct {
	Health *HealthHandler
}

func New() *Handler {
	return &Handler{
		Health: NewHealthHandler(),
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", h.Health.Health)

	// API v1 routes
	v1 := e.Group("/api/v1")

	// Auth routes
	auth := v1.Group("/auth")
	_ = auth // TODO: register auth routes

	// User routes
	users := v1.Group("/users")
	_ = users // TODO: register user routes

	// Client routes (admin)
	clients := v1.Group("/clients")
	_ = clients // TODO: register client routes

	// OAuth routes
	oauth := e.Group("/oauth")
	_ = oauth // TODO: register oauth routes
}
