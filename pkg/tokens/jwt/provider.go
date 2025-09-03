package jwt

import (
	"errors"
	"fmt"
	"sso-server/pkg/tokens"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Provider struct {
	config tokens.JWTConfig
}

func NewProvider(config tokens.JWTConfig) *Provider {
	return &Provider{
		config: config,
	}
}

func (p *Provider) GenerateAccessToken(req tokens.CreateTokenRequest) (string, tokens.TokenClaims, error) {
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

func (p *Provider) signToken(claims *CustomJwtClaims) (string, error) {
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

	return token.SignedString(p.config.SecretKey)
}

func (p *Provider) getKeyID() string {
	// This is a placeholder for key ID generation
	// In a production environment, you might want to implement proper key rotation
	// and return the current key ID
	return ""
}

func (p *Provider) GenerateRefreshTokenFromAccessToken(accessTokenClaims tokens.TokenClaims) (string, tokens.TokenClaims, error) {
	tokenID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(p.config.RefreshTokenExpiry)

	// Generate unique refresh secret for token family management
	refreshSecret := uuid.New().String()

	// Create refresh token with MINIMAL scopes - this is the key difference!
	refreshClaims := &CustomJwtClaims{
		UserID:        accessTokenClaims.GetUserID(),
		ApplicationID: accessTokenClaims.GetApplicationID(),
		Email:         accessTokenClaims.GetEmail(),

		Scopes:        []string{"refresh"},
		TokenType:     tokens.RefreshTokenType,
		TokenID:       tokenID,
		AccessTokenID: accessTokenClaims.GetTokenID(),
		RefreshSecret: refreshSecret,

		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   accessTokenClaims.GetUserID().String(),
			Audience:  jwt.ClaimStrings{"auth-server"}, // Only auth server, not resource servers
			Issuer:    p.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	tokenString, err := p.signToken(refreshClaims)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, refreshClaims, nil
}

func (p *Provider) ValidateAccessToken(tokenString string) (tokens.TokenClaims, error) {
	claims, err := p.parseAndValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != tokens.AccessTokenType {
		return nil, fmt.Errorf("%w: expected access token, got %s", tokens.ErrInvalidTokenType, claims.TokenType)
	}

	return claims, nil
}

func (p *Provider) ValidateRefreshToken(tokenString string) (tokens.TokenClaims, error) {
	claims, err := p.parseAndValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != tokens.RefreshTokenType {
		return nil, fmt.Errorf("%w: expected refresh token, got %s", tokens.ErrInvalidTokenType, claims.TokenType)
	}

	return claims, nil
}

func (p *Provider) parseAndValidateToken(tokenString string) (*CustomJwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		expectedAlg := p.config.Algorithm
		if expectedAlg == "" {
			expectedAlg = "HS256"
		}

		if token.Method.Alg() != expectedAlg {
			return nil, fmt.Errorf("%w: unexpected signing method %v", tokens.ErrInvalidSignature, token.Header["alg"])
		}

		return p.config.SecretKey, nil
	})

	if err != nil {
		// Check for specific JWT v5 error types using errors.Is
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, tokens.ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, tokens.ErrTokenInvalid
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, tokens.ErrTokenInvalid
		case errors.Is(err, jwt.ErrTokenUnverifiable):
			return nil, tokens.ErrTokenInvalid
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, tokens.ErrInvalidSignature
		case errors.Is(err, jwt.ErrTokenRequiredClaimMissing):
			return nil, tokens.ErrInvalidClaims
		case errors.Is(err, jwt.ErrTokenInvalidAudience):
			return nil, tokens.ErrInvalidClaims
		case errors.Is(err, jwt.ErrTokenInvalidIssuer):
			return nil, tokens.ErrInvalidClaims
		case errors.Is(err, jwt.ErrTokenInvalidSubject):
			return nil, tokens.ErrInvalidClaims
		case errors.Is(err, jwt.ErrTokenInvalidId):
			return nil, tokens.ErrInvalidClaims
		case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
			return nil, tokens.ErrTokenInvalid
		default:
			return nil, fmt.Errorf("%w: %v", tokens.ErrTokenInvalid, err)
		}
	}

	if !token.Valid {
		return nil, tokens.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*CustomJwtClaims)
	if !ok {
		return nil, fmt.Errorf("%w: failed to extract claims", tokens.ErrInvalidClaims)
	}

	// Additional validation for issuer if configured
	if p.config.Issuer != "" && claims.Issuer != p.config.Issuer {
		return nil, fmt.Errorf("%w: issuer mismatch", tokens.ErrInvalidClaims)
	}

	return claims, nil
}

func (p *Provider) GetTokenInfo(tokenString string) (tokens.TokenValidationResult, error) {
	claims, err := p.parseAndValidateToken(tokenString)

	result := tokens.TokenValidationResult{
		Valid:  err == nil,
		Claims: nil,
		Error:  err,
	}

	if err == nil {
		result.Claims = claims
	}

	return result, nil
}

func (p *Provider) ExtractClaimsWithoutValidation(tokenString string) (tokens.TokenClaims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, _, err := parser.ParseUnverified(tokenString, &CustomJwtClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomJwtClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}

	return claims, nil
}

func (p *Provider) GetProviderType() tokens.TokenProviderType {
	return tokens.TokenProviderJWT
}
