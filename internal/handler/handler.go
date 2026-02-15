package handler

import "github.com/labstack/echo/v4"

type Handler struct {
	Health *HealthHandler
	Auth   *AuthHandler
	User   *UserHandler
	Client *ClientHandler
	OAuth  *OAuthHandler
}

func New() *Handler {
	return &Handler{
		Health: NewHealthHandler(),
		Auth:   NewAuthHandler(),
		User:   NewUserHandler(),
		Client: NewClientHandler(),
		OAuth:  NewOAuthHandler(),
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	// Health check
	e.GET("/health", h.Health.Health)

	// API v1 routes
	v1 := e.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.POST("/register", h.Auth.Register)
	auth.POST("/login", h.Auth.Login)
	auth.POST("/refresh", h.Auth.Refresh)
	auth.POST("/logout", h.Auth.Logout) // TODO: add auth middleware

	// User routes (protected)
	users := v1.Group("/users")
	// TODO: add auth middleware
	users.GET("/me", h.User.GetMe)
	users.PATCH("/me", h.User.UpdateMe)
	users.PUT("/me/password", h.User.ChangePassword)

	// Client routes (admin protected)
	clients := v1.Group("/clients")
	// TODO: add admin auth middleware
	clients.POST("", h.Client.Create)
	clients.GET("", h.Client.List)
	clients.GET("/:id", h.Client.Get)
	clients.DELETE("/:id", h.Client.Delete)

	// OAuth routes
	oauth := e.Group("/oauth")
	oauth.GET("/authorize", h.OAuth.Authorize)
	oauth.POST("/token", h.OAuth.Token)
	oauth.POST("/revoke", h.OAuth.Revoke)
	oauth.GET("/userinfo", h.OAuth.UserInfo) // TODO: add auth middleware
}
