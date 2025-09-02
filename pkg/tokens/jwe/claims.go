package jwe

import (
	"sso-server/pkg/tokens"
	"time"

	"github.com/google/uuid"
)

type CustomJweClaims struct {
	UserID        uuid.UUID
	ApplicationID uuid.UUID
	Email         string
	Scopes        []string
	TokenType     tokens.TokenType
	TokenID       uuid.UUID
	AccessTokenID *uuid.UUID
	RefreshSecret string
	CustomClaims  map[string]interface{}

	// Standard claims
	Subject   string
	Issuer    string
	ID        string
	IssuedAt  time.Time
	ExpiresAt time.Time
	NotBefore time.Time
	Audience  []string
}

func (c *CustomJweClaims) GetUserID() uuid.UUID {
	return c.UserID
}

func (c *CustomJweClaims) GetApplicationID() uuid.UUID {
	return c.ApplicationID
}

func (c *CustomJweClaims) GetEmail() string {
	return c.Email
}

func (c *CustomJweClaims) GetScopes() []string {
	if c.Scopes == nil {
		return []string{}
	}
	return c.Scopes
}

func (c *CustomJweClaims) GetTokenType() tokens.TokenType {
	return c.TokenType
}

func (c *CustomJweClaims) GetTokenID() uuid.UUID {
	return c.TokenID
}

func (c *CustomJweClaims) GeTokenExpirationTime() time.Time {
	return c.ExpiresAt
}

func (c *CustomJweClaims) GetTokenIssuedAt() time.Time {
	return c.IssuedAt
}

func (c *CustomJweClaims) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

func (c *CustomJweClaims) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"user_id":        c.UserID.String(),
		"application_id": c.ApplicationID.String(),
		"email":          c.Email,
		"scopes":         c.Scopes,
		"token_type":     string(c.TokenType),
		"token_id":       c.TokenID.String(),
		"sub":            c.Subject,
		"iss":            c.Issuer,
		"jti":            c.ID,
		"iat":            c.IssuedAt.Unix(),
		"exp":            c.ExpiresAt.Unix(),
		"nbf":            c.NotBefore.Unix(),
		"aud":            c.Audience,
	}

	if c.AccessTokenID != nil {
		result["access_token_id"] = c.AccessTokenID.String()
	}
	if c.RefreshSecret != "" {
		result["refresh_secret"] = c.RefreshSecret
	}

	if c.CustomClaims != nil {
		for k, v := range c.CustomClaims {
			result[k] = v
		}
	}

	return result
}
