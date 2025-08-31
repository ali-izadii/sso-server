package repositories

import (
	"context"
	"sso-server/internal/domain/models"

	"github.com/google/uuid"
)

type TokenRepository interface {
	CreateAccessToken(ctx context.Context, token *models.AccessToken) error
	GetAccessToken(ctx context.Context, token string) (*models.AccessToken, error)
	GetAccessTokenByID(ctx context.Context, id uuid.UUID) (*models.AccessToken, error)
	UpdateAccessToken(ctx context.Context, token *models.AccessToken) error
	RevokeAccessToken(ctx context.Context, token string) error
	RevokeAccessTokenByID(ctx context.Context, id uuid.UUID) error

	//CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	//GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	//GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (*models.RefreshToken, error)
	//UpdateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	//RevokeRefreshToken(ctx context.Context, token string) error
	//RevokeRefreshTokenByID(ctx context.Context, id uuid.UUID) error
	//
	//GetAccessTokensByUserID(ctx context.Context, userID uuid.UUID) ([]*models.AccessToken, error)
	//GetRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) ([]*models.RefreshToken, error)
	//GetTokensByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.AccessToken, []*models.RefreshToken, error)
	//
	//RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	//RevokeAllApplicationTokens(ctx context.Context, applicationID uuid.UUID) error
	//RevokeTokenPair(ctx context.Context, accessTokenID uuid.UUID) error // Revokes both access and refresh tokens
	//
	//DeleteExpiredTokens(ctx context.Context) (int64, error) // Returns number of deleted tokens
	//DeleteRevokedTokens(ctx context.Context, olderThan time.Time) (int64, error)
	//
	//ValidateAccessToken(ctx context.Context, token string) (*models.AccessToken, error)
	//ValidateRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	//
	//CountActiveTokens(ctx context.Context) (int64, error)
	//CountTokensByUser(ctx context.Context, userID uuid.UUID) (int64, int64, error) // access, refresh counts
}
