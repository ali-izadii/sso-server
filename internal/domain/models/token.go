package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	AccessTokenType  TokenType = "access_token"
	RefreshTokenType TokenType = "refresh_token"
)

type AccessToken struct {
	ID            uuid.UUID `db:"id"`
	Token         string    `db:"tokens"`
	UserID        uuid.UUID `db:"user_id"`
	ApplicationID uuid.UUID `db:"application_id"`
	Scopes        string    `db:"scopes"`
	ExpiresAt     time.Time `db:"expires_at"`
	Revoked       bool      `db:"revoked"`
	CreatedAt     time.Time `db:"created_at"`
}

func (at *AccessToken) IsExpired() bool {
	return time.Now().After(at.ExpiresAt)
}

func (at *AccessToken) IsValid() bool {
	return !at.IsExpired() && !at.Revoked
}

func (at *AccessToken) ScopesAsSlice() []string {
	if at.Scopes == "" {
		return []string{}
	}
	return splitScopes(at.Scopes)
}

func splitScopes(scopes string) []string {
	if scopes == "" {
		return []string{}
	}
	var result []string
	for _, scope := range strings.Fields(scopes) {
		if scope != "" {
			result = append(result, scope)
		}
	}
	return result
}

type RefreshToken struct {
	ID            uuid.UUID  `db:"id"`
	Token         string     `db:"tokens"`
	UserID        uuid.UUID  `db:"user_id"`
	ApplicationID uuid.UUID  `db:"application_id"`
	AccessTokenID *uuid.UUID `db:"access_token_id"`
	ExpiresAt     time.Time  `db:"expires_at"`
	Revoked       bool       `db:"revoked"`
	CreatedAt     time.Time  `db:"created_at"`
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.Revoked
}

type TokenResponse struct {
	Token     string
	TokenType string
	ExpiresAt time.Time
	ExpiresIn int64
}

type TokenPair struct {
	AccessToken  TokenResponse
	RefreshToken TokenResponse
}

type OAuthTokenRequest struct {
	GrantType         string
	AuthorizationCode string
	RedirectURI       string
	RefreshToken      string
	ClientID          string
	ClientSecret      string
	Scope             string
}

type OAuthTokenResponse struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
	Scope        string
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RevokeTokenRequest struct {
	Token     string `json:"tokens" binding:"required"`
	TokenType string `json:"token_type,omitempty"` // "access" or "refresh"
}

func ScopesAsString(scopes []string) string {
	if len(scopes) == 0 {
		return "openid profile email"
	}
	return joinScopes(scopes)
}

func joinScopes(scopes []string) string {
	return strings.Join(scopes, " ")
}
