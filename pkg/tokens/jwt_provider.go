package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtProvider struct {
	secretKey          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	issuer             string
	signingMethod      jwt.SigningMethod
}

func NewJWTProvider(config JWTConfig) (TokenProvider, error) {
	if config.SecretKey == "" {
		return nil, errors.New("JWT secret key cannot be empty")
	}

	// Set defaults
	if config.AccessTokenExpiry == 0 {
		config.AccessTokenExpiry = 15 * time.Minute
	}
	if config.RefreshTokenExpiry == 0 {
		config.RefreshTokenExpiry = 7 * 24 * time.Hour
	}
	if config.Issuer == "" {
		config.Issuer = "sso-server"
	}
	if config.Algorithm == "" {
		config.Algorithm = "HS256"
	}

	// Determine signing method
	var signingMethod jwt.SigningMethod
	switch config.Algorithm {
	case "HS256":
		signingMethod = jwt.SigningMethodHS256
	case "HS384":
		signingMethod = jwt.SigningMethodHS384
	case "HS512":
		signingMethod = jwt.SigningMethodHS512
	default:
		return nil, fmt.Errorf("unsupported signing algorithm: %s", config.Algorithm)
	}

	return &jwtProvider{
		secretKey:          []byte(config.SecretKey),
		accessTokenExpiry:  config.AccessTokenExpiry,
		refreshTokenExpiry: config.RefreshTokenExpiry,
		issuer:             config.Issuer,
		signingMethod:      signingMethod,
	}, nil
}

func (j jwtProvider) GenerateAccessToken(req CreateTokenRequest) (string, *TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (j jwtProvider) GenerateRefreshToken(req CreateTokenRequest, accessTokenID uuid.UUID) (string, *TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (j jwtProvider) ValidateToken(tokenString string) (*TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (j jwtProvider) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (j jwtProvider) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (j jwtProvider) GetTokenInfo(tokenString string) (*TokenValidationResult, error) {
	//TODO implement me
	panic("implement me")
}

func (j jwtProvider) ExtractClaimsWithoutValidation(tokenString string) (*TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (j jwtProvider) GetTokenExpiry(tokenType TokenType) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (j jwtProvider) GetProviderType() TokenProviderType {
	//TODO implement me
	panic("implement me")
}

type JWTClaims struct {
	UserID        uuid.UUID
	Email         string
	ApplicationID uuid.UUID
	Scopes        []string
	TokenType     TokenType
	TokenID       uuid.UUID
	jwt.RegisteredClaims
}

func (claim JWTClaims) GetUserID() uuid.UUID {
	return claim.UserID
}

func (claim JWTClaims) GetApplicationID() uuid.UUID {
	return claim.ApplicationID
}

func (claim JWTClaims) GetEmail() string {
	return claim.Email
}

func (claim JWTClaims) GetScopes() []string {
	return claim.Scopes
}

func (claim JWTClaims) GetTokenType() TokenType {
	return claim.TokenType
}

func (claim JWTClaims) GetTokenID() uuid.UUID {
	return claim.TokenID
}

func (claim JWTClaims) GetExpiresAt() time.Time {
	if claim.ExpiresAt == nil {
		return time.Time{}
	}
	return claim.ExpiresAt.Time
}

func (claim JWTClaims) GetIssuedAt() time.Time {
	if claim.IssuedAt == nil {
		return time.Time{}
	}
	return claim.ExpiresAt.Time
}

func (claim JWTClaims) IsExpired() bool {
	return time.Now().After(claim.GetExpiresAt())
}

func (claim JWTClaims) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id":        claim.GetUserID().String(),
		"application_id": claim.GetApplicationID().String(),
		"email":          claim.GetEmail(),
		"scopes":         claim.GetScopes(),
		"token_type":     string(claim.GetTokenType()),
		"token_id":       claim.GetTokenID().String(),
		"expires_at":     claim.GetExpiresAt(),
		"issued_at":      claim.GetIssuedAt(),
		"issuer":         claim.Issuer,
		"subject":        claim.Subject,
		"audience":       claim.Audience,
	}
}
