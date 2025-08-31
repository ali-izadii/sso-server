package models

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessTokenType  TokenType = "access_token"
	RefreshTokenType TokenType = "refresh_token"
)

type AccessToken struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Token         string    `json:"token" db:"token"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	ApplicationID uuid.UUID `json:"application_id" db:"application_id"`
	Scopes        string    `json:"scopes" db:"scopes"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	Revoked       bool      `json:"revoked" db:"revoked"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
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
	ID            uuid.UUID  `json:"id" db:"id"`
	Token         string     `json:"token" db:"token"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	ApplicationID uuid.UUID  `json:"application_id" db:"application_id"`
	AccessTokenID *uuid.UUID `json:"access_token_id" db:"access_token_id"`
	ExpiresAt     time.Time  `json:"expires_at" db:"expires_at"`
	Revoked       bool       `json:"revoked" db:"revoked"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.Revoked
}

type JWTClaims struct {
	UserID        uuid.UUID `json:"user_id"`
	Email         string    `json:"email"`
	ApplicationID uuid.UUID `json:"application_id"`
	Scopes        []string  `json:"scopes"`
	TokenType     TokenType `json:"token_type"`
	TokenID       uuid.UUID `json:"token_id"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	ExpiresIn int64     `json:"expires_in"` // seconds until expiration
}

type TokenPair struct {
	AccessToken  TokenResponse `json:"access_token"`
	RefreshToken TokenResponse `json:"refresh_token"`
}

type CreateTokenRequest struct {
	UserID        uuid.UUID
	ApplicationID uuid.UUID
	Scopes        []string
	Email         string
}

type OAuthTokenRequest struct {
	GrantType         string `json:"grant_type" form:"grant_type" binding:"required"`
	AuthorizationCode string `json:"code" form:"code"`
	RedirectURI       string `json:"redirect_uri" form:"redirect_uri"`
	RefreshToken      string `json:"refresh_token" form:"refresh_token"`
	ClientID          string `json:"client_id" form:"client_id" binding:"required"`
	ClientSecret      string `json:"client_secret" form:"client_secret"`
	Scope             string `json:"scope" form:"scope"`
}

type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope,omitempty"`
}

type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
type TokenValidationResult struct {
	Valid         bool       `json:"valid"`
	Claims        *JWTClaims `json:"claims,omitempty"`
	Error         string     `json:"error,omitempty"`
	ExpiresAt     time.Time  `json:"expires_at,omitempty"`
	TokenType     TokenType  `json:"token_type,omitempty"`
	UserID        uuid.UUID  `json:"user_id,omitempty"`
	ApplicationID uuid.UUID  `json:"application_id,omitempty"`
	Scopes        []string   `json:"scopes,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RevokeTokenRequest struct {
	Token     string `json:"token" binding:"required"`
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
