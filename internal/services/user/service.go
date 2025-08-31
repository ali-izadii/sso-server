package user

import (
	"context"
	"fmt"
	"strings"

	"sso-server/internal/domain/models"
	"sso-server/internal/domain/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req models.CreateUserRequest) (*models.UserResponse, error)
	ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error
	VerifyEmail(ctx context.Context, id uuid.UUID) error
	DeactivateUser(ctx context.Context, id uuid.UUID) error
	ReactivateUser(ctx context.Context, id uuid.UUID) error
	ValidateCredentials(ctx context.Context, email, password string) (*models.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*models.UserResponse, int64, error)
}

type service struct {
	userRepo repositories.UserRepository
}

func NewService(userRepo repositories.UserRepository) Service {
	return &service{
		userRepo: userRepo,
	}
}

// RegisterUser creates a new user account
func (s *service) RegisterUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check if user already exists
	exists, err := s.userRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if email exists: %w", err)
	}

	if exists {
		return nil, models.ErrUserAlreadyExists
	}

	// Hash password
	passwordHash, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user parameters
	params := models.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	// Create user
	user, err := s.userRepo.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByEmail retrieves a user by email (returns full user for internal use)
func (s *service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	return s.userRepo.GetByEmail(ctx, email)
}

// UpdateUser updates user information
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, req models.CreateUserRequest) (*models.UserResponse, error) {
	// Get existing user
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	updatedUser := &models.User{
		ID:        existingUser.ID,
		Email:     existingUser.Email, // Email cannot be changed via this method
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Update user
	user, err := s.userRepo.Update(ctx, id, updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword changes user's password
func (s *service) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify old password
	if !s.verifyPassword(oldPassword, user.PasswordHash) {
		return models.ErrInvalidPassword
	}

	// Validate new password
	if len(newPassword) < 8 {
		return models.ErrPasswordTooShort
	}

	// Hash new password
	newPasswordHash, err := s.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	return s.userRepo.UpdatePassword(ctx, id, newPasswordHash)
}

// VerifyEmail marks user's email as verified
func (s *service) VerifyEmail(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.VerifyEmail(ctx, id)
}

// DeactivateUser deactivates a user account
func (s *service) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeactivateUser(ctx, id)
}

// ReactivateUser reactivates a user account
func (s *service) ReactivateUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.ReactivateUser(ctx, id)
}

// ValidateCredentials validates user credentials for authentication
func (s *service) ValidateCredentials(ctx context.Context, email, password string) (*models.User, error) {
	// Get user by email
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		if err == models.ErrUserNotFound {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, models.ErrUserInactive
	}

	// Verify password
	if !s.verifyPassword(password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return user, nil
}

// ListUsers retrieves users with pagination
func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]*models.UserResponse, int64, error) {
	// Get users
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Convert to response format
	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		response := user.ToResponse()
		responses[i] = &response
	}

	return responses, total, nil
}

// hashPassword hashes a password using bcrypt
func (s *service) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// verifyPassword verifies a password against its hash
func (s *service) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
