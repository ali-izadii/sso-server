package database

import (
	"context"
	"database/sql"
	"fmt"
	"sso-server/internal/domain/models"
	"sso-server/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) repositories.TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (r *tokenRepository) CreateAccessToken(ctx context.Context, token *models.AccessToken) error {
	query := `
		INSERT INTO access_tokens (id, tokens, user_id, application_id, scopes, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		token.ID,
		token.Token,
		token.UserID,
		token.ApplicationID,
		token.Scopes,
		token.ExpiresAt,
		token.Revoked,
		token.CreatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("access tokens already exists: %w", err)
		}
		return fmt.Errorf("failed to create access tokens: %w", err)
	}

	return nil
}

func (r *tokenRepository) GetAccessToken(ctx context.Context, token string) (*models.AccessToken, error) {
	query := `
		SELECT id, tokens, user_id, application_id, scopes, expires_at, revoked, created_at
		FROM access_tokens
		WHERE tokens = $1
	`

	accessToken := &models.AccessToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&accessToken.ID,
		&accessToken.Token,
		&accessToken.UserID,
		&accessToken.ApplicationID,
		&accessToken.Scopes,
		&accessToken.ExpiresAt,
		&accessToken.Revoked,
		&accessToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrTokenNotFound
		}
		return nil, fmt.Errorf("failed to get access tokens: %w", err)
	}

	return accessToken, nil
}

func (r *tokenRepository) GetAccessTokenByID(ctx context.Context, id uuid.UUID) (*models.AccessToken, error) {
	query := `
		SELECT id, tokens, user_id, application_id, scopes, expires_at, revoked, created_at
		FROM access_tokens
		WHERE id = $1
	`

	accessToken := &models.AccessToken{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&accessToken.ID,
		&accessToken.Token,
		&accessToken.UserID,
		&accessToken.ApplicationID,
		&accessToken.Scopes,
		&accessToken.ExpiresAt,
		&accessToken.Revoked,
		&accessToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrTokenNotFound
		}
		return nil, fmt.Errorf("failed to get access tokens by ID: %w", err)
	}

	return accessToken, nil
}

func (r *tokenRepository) UpdateAccessToken(ctx context.Context, token *models.AccessToken) error {
	query := `
		UPDATE access_tokens 
		SET scopes = $2, expires_at = $3, revoked = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		token.ID,
		token.Scopes,
		token.ExpiresAt,
		token.Revoked,
	)

	if err != nil {
		return fmt.Errorf("failed to update access tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrTokenNotFound
	}

	return nil
}

func (r *tokenRepository) RevokeAccessToken(ctx context.Context, token string) error {
	query := `
		UPDATE access_tokens 
		SET revoked = true
		WHERE tokens = $1
	`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to revoke access tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrTokenNotFound
	}

	return nil
}

func (r *tokenRepository) RevokeAccessTokenByID(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE access_tokens 
		SET revoked = true
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke access tokens by ID: %w", err)
	}
	return nil
}
