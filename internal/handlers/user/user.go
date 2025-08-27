package user

import (
	"net/http"
	"strconv"

	"sso-server/internal/domain/models"
	"sso-server/internal/services/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	userService user.Service
}

// NewHandler creates a new user handler
func NewHandler(userService user.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

// RegisterRoutes registers user routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", h.Register)
		userRoutes.GET("/profile/:id", h.GetProfile)
		userRoutes.PUT("/profile/:id", h.UpdateProfile)
		userRoutes.POST("/change-password", h.ChangePassword)
		userRoutes.POST("/verify-email/:id", h.VerifyEmail)
		userRoutes.GET("/", h.ListUsers)
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.RegisterUser(c.Request.Context(), req)
	if err != nil {
		switch err {
		case models.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
		case models.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid email address",
			})
		case models.ErrPasswordTooShort:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must be at least 8 characters long",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to register user",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// GetProfile retrieves user profile by ID
func (h *Handler) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve user profile",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateProfile updates user profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update user profile",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// ChangePasswordRequest represents the request to change password
type ChangePasswordRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword handles password change
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	id, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	err = h.userService.ChangePassword(c.Request.Context(), id, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		case models.ErrInvalidPassword:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Current password is incorrect",
			})
		case models.ErrPasswordTooShort:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "New password must be at least 8 characters long",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to change password",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// VerifyEmail handles email verification
func (h *Handler) VerifyEmail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	err = h.userService.VerifyEmail(c.Request.Context(), id)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to verify email",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users []*models.UserResponse `json:"users"`
	Total int64                  `json:"total"`
	Limit int                    `json:"limit"`
	Page  int                    `json:"page"`
}

// ListUsers handles user listing with pagination
func (h *Handler) ListUsers(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Cap at 100
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	users, total, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	response := ListUsersResponse{
		Users: users,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	c.JSON(http.StatusOK, response)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
