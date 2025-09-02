package jwt

import (
	"context"
	"fmt"
	"sso-server/pkg/tokens"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Provider struct {
	config tokens.JWTConfig
}

func (p Provider) GenerateAccessToken(ctx context.Context, req tokens.CreateTokenRequest) (string, tokens.TokenClaims, error) {
	tokenID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(p.config.AccessTokenExpiry)

	claims := &CustomJwtClaims{
		UserID:        req.UserID,
		ApplicationID: req.ApplicationID,
		Email:         req.Email,
		Scopes:        req.Scopes,
		TokenType:     tokens.AccessTokenType,
		TokenID:       tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   req.UserID.String(),
			Audience:  jwt.ClaimStrings{req.ApplicationID.String()},
			Issuer:    p.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	tokenString, err := p.signToken(claims)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, claims, nil
}

func (p Provider) signToken(claims *CustomJwtClaims) (string, error) {
	var method jwt.SigningMethod
	switch p.config.Algorithm {
	case "HS256", "":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	default:
		return "", fmt.Errorf("unsupported signing algorithm: %s", p.config.Algorithm)
	}

	token := jwt.NewWithClaims(method, claims)

	if keyID := p.getKeyID(); keyID != "" {
		token.Header["kid"] = keyID
	}

	return token.SignedString([]byte(p.config.SecretKey))
}

func (p Provider) getKeyID() string {
	// This is a placeholder for key ID generation
	// In a production environment, you might want to implement proper key rotation
	// and return the current key ID
	return ""
}

func (p Provider) GenerateRefreshToken(ctx context.Context, req tokens.CreateTokenRequest, accessTokenID uuid.UUID) (string, tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) ValidateToken(ctx context.Context, tokenString string) (tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) ValidateAccessToken(ctx context.Context, tokenString string) (tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) ValidateRefreshToken(ctx context.Context, tokenString string) (tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) GetTokenInfo(ctx context.Context, tokenString string) (tokens.TokenValidationResult, error) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) ExtractClaimsWithoutValidation(tokenString string) (tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) GetTokenExpiry(tokenType tokens.TokenType) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (p Provider) GetProviderType() tokens.TokenProviderType {
	//TODO implement me
	panic("implement me")
}
