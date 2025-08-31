package token

import (
	"context"
	"fmt"
	"sso-server/internal/domain/models"
	"sso-server/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateTokenPair(ctx context.Context, req models.CreateTokenRequest) (*models.TokenPair, error)

	//ValidateToken(ctx context.Context, tokenString string) (*models.TokenValidationResult, error)
	//ValidateAccessToken(ctx context.Context, tokenString string) (*models.TokenValidationResult, error)
	//ValidateRefreshToken(ctx context.Context, tokenString string) (*models.TokenValidationResult, error)
	//
	//RefreshTokens(ctx context.Context, refreshTokenString string) (*models.TokenPair, error)
	//
	//RevokeToken(ctx context.Context, tokenString string) error
	//RevokeTokenPair(ctx context.Context, accessTokenID uuid.UUID) error
	//RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	//RevokeAllApplicationTokens(ctx context.Context, applicationID uuid.UUID) error
	//
	//GetTokenInfo(ctx context.Context, tokenString string) (*models.TokenValidationResult, error)
	//GetUserTokens(ctx context.Context, userID uuid.UUID) ([]*models.AccessToken, []*models.RefreshToken, error)
	//
	//CleanupExpiredTokens(ctx context.Context) (int64, error)
	//CleanupRevokedTokens(ctx context.Context, olderThan time.Time) (int64, error)
	//
	//GetTokenStats(ctx context.Context) (*TokenStats, error)
	//GetUserTokenStats(ctx context.Context, userID uuid.UUID) (*UserTokenStats, error)
}

//type TokenStats struct {
//	ActiveTokens  int64 `json:"active_tokens"`
//	ExpiredTokens int64 `json:"expired_tokens"`
//	RevokedTokens int64 `json:"revoked_tokens"`
//	TotalTokens   int64 `json:"total_tokens"`
//}
//
//type UserTokenStats struct {
//	UserID        uuid.UUID `json:"user_id"`
//	ActiveAccess  int64     `json:"active_access"`
//	ActiveRefresh int64     `json:"active_refresh"`
//	TotalActive   int64     `json:"total_active"`
//}

type service struct {
	tokenRepo        repositories.TokenRepository
	jwtService       *jwt.Service
	maxTokensPerUser int
	cleanupInterval  time.Duration
}

func NewService(tokenRepo repositories.TokenRepository, jwtService *jwt.Service) Service {
	return &service{
		tokenRepo:        tokenRepo,
		jwtService:       jwtService,
		maxTokensPerUser: 10, // Default limit
		cleanupInterval:  24 * time.Hour,
	}
}

func (s *service) CreateTokenPair(ctx context.Context, req models.CreateTokenRequest) (*models.TokenPair, error) {
	// Check tokens limits
	accessCount, refreshCount, err := s.tokenRepo.CountTokensByUser(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to count user tokens: %w", err)
	}

	if accessCount >= int64(s.maxTokensPerUser) {
		return nil, models.ErrTooManyTokens
	}

	// Generate JWT access tokens
	accessTokenString, accessClaims, err := s.jwtService.GenerateAccessToken(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access tokens: %w", err)
	}

	// Create access tokens record
	accessToken := &models.AccessToken{
		ID:            accessClaims.TokenID,
		Token:         accessTokenString,
		UserID:        req.UserID,
		ApplicationID: req.ApplicationID,
		Scopes:        models.ScopesAsString(req.Scopes),
		ExpiresAt:     accessClaims.ExpiresAt.Time,
		Revoked:       false,
		CreatedAt:     time.Now(),
	}

	if err := s.tokenRepo.CreateAccessToken(ctx, accessToken); err != nil {
		return nil, fmt.Errorf("failed to store access tokens: %w", err)
	}

	// Generate JWT refresh tokens
	refreshTokenString, refreshClaims, err := s.jwtService.GenerateRefreshToken(req, accessToken.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh tokens: %w", err)
	}

	// Create refresh tokens record
	refreshToken := &models.RefreshToken{
		ID:            refreshClaims.TokenID,
		Token:         refreshTokenString,
		UserID:        req.UserID,
		ApplicationID: req.ApplicationID,
		AccessTokenID: &accessToken.ID,
		ExpiresAt:     refreshClaims.ExpiresAt.Time,
		Revoked:       false,
		CreatedAt:     time.Now(),
	}

	if err := s.tokenRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		// Clean up access tokens if refresh tokens creation fails
		_ = s.tokenRepo.RevokeAccessTokenByID(ctx, accessToken.ID)
		return nil, fmt.Errorf("failed to store refresh tokens: %w", err)
	}

	// Build response
	tokenPair := &models.TokenPair{
		AccessToken: models.TokenResponse{
			Token:     accessTokenString,
			TokenType: "Bearer",
			ExpiresAt: accessToken.ExpiresAt,
			ExpiresIn: int64(time.Until(accessToken.ExpiresAt).Seconds()),
		},
		RefreshToken: models.TokenResponse{
			Token:     refreshTokenString,
			TokenType: "Bearer",
			ExpiresAt: refreshToken.ExpiresAt,
			ExpiresIn: int64(time.Until(refreshToken.ExpiresAt).Seconds()),
		},
	}

	return tokenPair, nil
}

// ValidateToken validates any tokens (access or refresh)
func (s *service) ValidateToken(ctx context.Context, tokenString string) (*models.TokenValidationResult, error) {
	// First try JWT validation
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return &models.TokenValidationResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	// Then check database status based on tokens type
	if claims.TokenType == models.TokenTypeAccess {
		return s.validateAccessTokenInDB(ctx, tokenString, claims)
	} else if claims.TokenType == models.TokenTypeRefresh {
		return s.validateRefreshTokenInDB(ctx, tokenString, claims)
	}

	return &models.TokenValidationResult{
		Valid: false,
		Error: "unknown tokens type",
	}, nil
}

// ValidateAccessToken validates an access tokens specifically
func (s *service) ValidateAccessToken(ctx context.Context, tokenString string) (*models.TokenValidationResult, error) {
	// Validate JWT signature and claims
	claims, err := s.jwtService.ValidateAccessToken(tokenString)
	if err != nil {
		return &models.TokenValidationResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return s.validateAccessTokenInDB(ctx, tokenString, claims)
}

// ValidateRefreshToken validates a refresh tokens specifically
func (s *service) ValidateRefreshToken(ctx context.Context, tokenString string) (*models.TokenValidationResult, error) {
	// Validate JWT signature and claims
	claims, err := s.jwtService.ValidateRefreshToken(tokenString)
	if err != nil {
		return &models.TokenValidationResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return s.validateRefreshTokenInDB(ctx, tokenString, claims)
}

// validateAccessTokenInDB validates access tokens against database
func (s *service) validateAccessTokenInDB(ctx context.Context, tokenString string, claims *models.JWTClaims) (*models.TokenValidationResult, error) {
	// Check database for tokens status
	dbToken, err := s.tokenRepo.ValidateAccessToken(ctx, tokenString)
	if err != nil {
		if err == models.ErrTokenInvalid {
			return &models.TokenValidationResult{
				Valid: false,
				Error: "tokens not found or invalid",
			}, nil
		}
		return nil, fmt.Errorf("database validation failed: %w", err)
	}

	// Token is valid
	return &models.TokenValidationResult{
		Valid:         true,
		Claims:        claims,
		ExpiresAt:     dbToken.ExpiresAt,
		TokenType:     models.TokenTypeAccess,
		UserID:        dbToken.UserID,
		ApplicationID: dbToken.ApplicationID,
		Scopes:        dbToken.ScopesAsSlice(),
	}, nil
}

// validateRefreshTokenInDB validates refresh tokens against database
func (s *service) validateRefreshTokenInDB(ctx context.Context, tokenString string, claims *models.JWTClaims) (*models.TokenValidationResult, error) {
	// Check database for tokens status
	dbToken, err := s.tokenRepo.ValidateRefreshToken(ctx, tokenString)
	if err != nil {
		if err == models.ErrTokenInvalid {
			return &models.TokenValidationResult{
				Valid: false,
				Error: "tokens not found or invalid",
			}, nil
		}
		return nil, fmt.Errorf("database validation failed: %w", err)
	}

	// Token is valid
	return &models.TokenValidationResult{
		Valid:         true,
		Claims:        claims,
		ExpiresAt:     dbToken.ExpiresAt,
		TokenType:     models.TokenTypeRefresh,
		UserID:        dbToken.UserID,
		ApplicationID: dbToken.ApplicationID,
		Scopes:        []string{}, // Refresh tokens don't have scopes directly
	}, nil
}

// RefreshTokens creates new tokens using a valid refresh tokens
func (s *service) RefreshTokens(ctx context.Context, refreshTokenString string) (*models.TokenPair, error) {
	// Validate refresh tokens
	validationResult, err := s.ValidateRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("refresh tokens validation failed: %w", err)
	}

	if !validationResult.Valid {
		return nil, fmt.Errorf("invalid refresh tokens: %s", validationResult.Error)
	}

	// Get the refresh tokens from database to access linked access tokens
	dbRefreshToken, err := s.tokenRepo.GetRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	// Get the original access tokens to inherit scopes
	var scopes []string
	if dbRefreshToken.AccessTokenID != nil {
		originalAccessToken, err := s.tokenRepo.GetAccessTokenByID(ctx, *dbRefreshToken.AccessTokenID)
		if err == nil {
			scopes = originalAccessToken.ScopesAsSlice()
		}
	}
	if len(scopes) == 0 {
		scopes = []string{"openid", "profile", "email"} // Default scopes
	}

	// Create new tokens pair
	req := models.CreateTokenRequest{
		UserID:        validationResult.UserID,
		ApplicationID: validationResult.ApplicationID,
		Scopes:        scopes,
		Email:         validationResult.Claims.Email,
	}

	tokenPair, err := s.CreateTokenPair(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create new tokens pair: %w", err)
	}

	// Revoke the old refresh tokens (and optionally the old access tokens)
	if err := s.tokenRepo.RevokeRefreshToken(ctx, refreshTokenString); err != nil {
		// Log warning but don't fail the operation
		fmt.Printf("Warning: failed to revoke old refresh tokens: %v\n", err)
	}

	// Optionally revoke the old access tokens
	if dbRefreshToken.AccessTokenID != nil {
		if err := s.tokenRepo.RevokeAccessTokenByID(ctx, *dbRefreshToken.AccessTokenID); err != nil {
			// Log warning but don't fail the operation
			fmt.Printf("Warning: failed to revoke old access tokens: %v\n", err)
		}
	}

	return tokenPair, nil
}

// RevokeToken revokes a tokens by its string value
func (s *service) RevokeToken(ctx context.Context, tokenString string) error {
	// Try to determine tokens type by parsing (without full validation)
	tokenInfo, err := s.jwtService.GetTokenInfo(tokenString)
	if err != nil {
		return fmt.Errorf("failed to get tokens info: %w", err)
	}

	if tokenInfo.Claims.TokenType == models.TokenTypeAccess {
		return s.tokenRepo.RevokeAccessToken(ctx, tokenString)
	} else if tokenInfo.Claims.TokenType == models.TokenTypeRefresh {
		return s.tokenRepo.RevokeRefreshToken(ctx, tokenString)
	}

	return models.ErrInvalidTokenType
}

// RevokeTokenPair revokes both access and refresh tokens for a tokens pair
func (s *service) RevokeTokenPair(ctx context.Context, accessTokenID uuid.UUID) error {
	return s.tokenRepo.RevokeTokenPair(ctx, accessTokenID)
}

// RevokeAllUserTokens revokes all tokens for a user
func (s *service) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return s.tokenRepo.RevokeAllUserTokens(ctx, userID)
}

// RevokeAllApplicationTokens revokes all tokens for an application
func (s *service) RevokeAllApplicationTokens(ctx context.Context, applicationID uuid.UUID) error {
	return s.tokenRepo.RevokeAllApplicationTokens(ctx, applicationID)
}

// GetTokenInfo returns tokens information without full validation
func (s *service) GetTokenInfo(ctx context.Context, tokenString string) (*models.TokenValidationResult, error) {
	return s.jwtService.GetTokenInfo(tokenString)
}

// GetUserTokens retrieves all tokens for a user
func (s *service) GetUserTokens(ctx context.Context, userID uuid.UUID) ([]*models.AccessToken, []*models.RefreshToken, error) {
	accessTokens, err := s.tokenRepo.GetAccessTokensByUserID(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user access tokens: %w", err)
	}

	refreshTokens, err := s.tokenRepo.GetRefreshTokensByUserID(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user refresh tokens: %w", err)
	}

	return accessTokens, refreshTokens, nil
}

// CleanupExpiredTokens removes expired tokens
func (s *service) CleanupExpiredTokens(ctx context.Context) (int64, error) {
	return s.tokenRepo.DeleteExpiredTokens(ctx)
}

// CleanupRevokedTokens removes old revoked tokens
func (s *service) CleanupRevokedTokens(ctx context.Context, olderThan time.Time) (int64, error) {
	return s.tokenRepo.DeleteRevokedTokens(ctx, olderThan)
}

// GetTokenStats returns overall tokens statistics
func (s *service) GetTokenStats(ctx context.Context) (*TokenStats, error) {
	activeTokens, err := s.tokenRepo.CountActiveTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count active tokens: %w", err)
	}

	return &TokenStats{
		ActiveTokens: activeTokens,
		// TODO: Add expired/revoked/total counts if needed
		TotalTokens: activeTokens,
	}, nil
}

// GetUserTokenStats returns tokens statistics for a specific user
func (s *service) GetUserTokenStats(ctx context.Context, userID uuid.UUID) (*UserTokenStats, error) {
	accessCount, refreshCount, err := s.tokenRepo.CountTokensByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count user tokens: %w", err)
	}

	return &UserTokenStats{
		UserID:        userID,
		ActiveAccess:  accessCount,
		ActiveRefresh: refreshCount,
		TotalActive:   accessCount + refreshCount,
	}, nil
}

// SetMaxTokensPerUser sets the maximum number of tokens per user
func (s *service) SetMaxTokensPerUser(max int) {
	s.maxTokensPerUser = max
}

// StartCleanupWorker starts a background worker to clean up expired tokens
func (s *service) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Clean up expired tokens
			deleted, err := s.CleanupExpiredTokens(ctx)
			if err != nil {
				fmt.Printf("Error cleaning up expired tokens: %v\n", err)
			} else if deleted > 0 {
				fmt.Printf("Cleaned up %d expired tokens\n", deleted)
			}

			// Clean up old revoked tokens (older than 30 days)
			olderThan := time.Now().Add(-30 * 24 * time.Hour)
			deletedRevoked, err := s.CleanupRevokedTokens(ctx, olderThan)
			if err != nil {
				fmt.Printf("Error cleaning up revoked tokens: %v\n", err)
			} else if deletedRevoked > 0 {
				fmt.Printf("Cleaned up %d old revoked tokens\n", deletedRevoked)
			}
		}
	}
}
