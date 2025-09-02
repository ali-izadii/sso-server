// pkg/tokens/jwt/claims.go
package jwt

import (
	"sso-server/pkg/tokens"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID        uuid.UUID
	ApplicationID uuid.UUID
	Email         string
	Scopes        []string
	TokenType     tokens.TokenType
	TokenID       uuid.UUID
	AccessTokenID *uuid.UUID
	RefreshSecret string
	CustomClaims  map[string]interface{}
	jwt.RegisteredClaims
}

func (c *JWTClaims) GetUserID() uuid.UUID {
	return c.UserID
}

func (c *JWTClaims) GetApplicationID() uuid.UUID {
	return c.ApplicationID
}

func (c *JWTClaims) GetEmail() string {
	return c.Email
}

func (c *JWTClaims) GetScopes() []string {
	if c.Scopes == nil {
		return []string{}
	}
	return c.Scopes
}

func (c *JWTClaims) GetTokenType() tokens.TokenType {
	return c.TokenType
}

func (c *JWTClaims) GetTokenID() uuid.UUID {
	return c.TokenID
}

func (c *JWTClaims) GetExpiresAt() time.Time {
	if c.ExpiresAt == nil {
		return time.Time{}
	}
	return c.ExpiresAt.Time
}

func (c *JWTClaims) GetIssuedAt() time.Time {
	if c.IssuedAt == nil {
		return time.Time{}
	}
	return c.IssuedAt.Time
}

func (c *JWTClaims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(c.ExpiresAt.Time)
}

func (c *JWTClaims) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"user_id":        c.UserID.String(),
		"application_id": c.ApplicationID.String(),
		"email":          c.Email,
		"scopes":         c.Scopes,
		"token_type":     string(c.TokenType),
		"token_id":       c.TokenID.String(),
	}

	// Add registered claims
	if c.Subject != "" {
		result["sub"] = c.Subject
	}
	if c.Issuer != "" {
		result["iss"] = c.Issuer
	}
	if c.ID != "" {
		result["jti"] = c.ID
	}
	if c.IssuedAt != nil {
		result["iat"] = c.IssuedAt.Unix()
	}
	if c.ExpiresAt != nil {
		result["exp"] = c.ExpiresAt.Unix()
	}
	if c.NotBefore != nil {
		result["nbf"] = c.NotBefore.Unix()
	}
	if len(c.Audience) > 0 {
		result["aud"] = c.Audience
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

func (c *JWTClaims) GetAccessTokenID() *uuid.UUID {
	return c.AccessTokenID
}

func (c *JWTClaims) GetRefreshSecret() string {
	return c.RefreshSecret
}

func (c *JWTClaims) GetCustomClaim(key string) interface{} {
	if c.CustomClaims == nil {
		return nil
	}
	return c.CustomClaims[key]
}

func (c *JWTClaims) SetCustomClaim(key string, value interface{}) {
	if c.CustomClaims == nil {
		c.CustomClaims = make(map[string]interface{})
	}
	c.CustomClaims[key] = value
}

func (c *JWTClaims) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func (c *JWTClaims) HasAllScopes(requiredScopes []string) bool {
	for _, required := range requiredScopes {
		if !c.HasScope(required) {
			return false
		}
	}
	return true
}

func (c *JWTClaims) HasAnyScope(requiredScopes []string) bool {
	for _, required := range requiredScopes {
		if c.HasScope(required) {
			return true
		}
	}
	return false
}

func (c *JWTClaims) Clone() *JWTClaims {
	clone := &JWTClaims{
		UserID:        c.UserID,
		ApplicationID: c.ApplicationID,
		Email:         c.Email,
		Scopes:        make([]string, len(c.Scopes)),
		TokenType:     c.TokenType,
		TokenID:       c.TokenID,
		RefreshSecret: c.RefreshSecret,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   c.Subject,
			Issuer:    c.Issuer,
			ID:        c.ID,
			IssuedAt:  c.IssuedAt,
			ExpiresAt: c.ExpiresAt,
			NotBefore: c.NotBefore,
		},
	}

	copy(clone.Scopes, c.Scopes)
	if len(c.Audience) > 0 {
		clone.Audience = make(jwt.ClaimStrings, len(c.Audience))
		copy(clone.Audience, c.Audience)
	}

	if c.AccessTokenID != nil {
		accessTokenID := *c.AccessTokenID
		clone.AccessTokenID = &accessTokenID
	}

	if c.CustomClaims != nil {
		clone.CustomClaims = make(map[string]interface{})
		for k, v := range c.CustomClaims {
			clone.CustomClaims[k] = v
		}
	}

	return clone
}
