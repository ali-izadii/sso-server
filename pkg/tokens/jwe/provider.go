package jwe

import (
	"fmt"
	"sso-server/pkg/tokens"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/google/uuid"
)

type Provider struct {
	config tokens.JWEConfig
}

func (p *Provider) GenerateAccessToken(req tokens.CreateTokenRequest) (string, tokens.TokenClaims, error) {
	tokenID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(p.config.AccessTokenExpiry)

	claims := &CustomJweClaims{
		UserID:        req.UserID,
		ApplicationID: req.ApplicationID,
		Email:         req.Email,
		Scopes:        req.Scopes,
		TokenType:     tokens.AccessTokenType,
		TokenID:       tokenID,
		Subject:       req.UserID.String(),
		Issuer:        p.config.Issuer,
		ID:            tokenID.String(),
		IssuedAt:      now,
		ExpiresAt:     expiresAt,
		NotBefore:     now,
		Audience:      []string{req.ApplicationID.String()},
	}

	tokenString, err := p.encryptToken(claims)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encrypt access token: %w", err)
	}

	return tokenString, claims, nil
}

func (p *Provider) encryptToken(claims *CustomJweClaims) (string, error) {
	// Marshal claims to JSON
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}

	// Create encrypter with the configured algorithms
	encrypter, err := jose.NewEncrypter(
		p.config.ContentEncryption,
		jose.Recipient{
			Algorithm: p.config.KeyEncryption,
			Key:       p.config.SecretKey,
		},
		(&jose.EncrypterOptions{}).WithType("JWE").WithContentType("JWT"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create encrypter: %w", err)
	}

	// Encrypt the payload
	object, err := encrypter.Encrypt(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt payload: %w", err)
	}

	// Serialize to compact format
	token, err := object.CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize encrypted token: %w", err)
	}

	return token, nil
}

func (p *Provider) GenerateRefreshToken(req tokens.CreateTokenRequest, accessTokenID uuid.UUID) (string, tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Provider) ValidateAccessToken(tokenString string) (tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Provider) ValidateRefreshToken(tokenString string) (tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Provider) GetTokenInfo(tokenString string) (tokens.TokenValidationResult, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Provider) ExtractClaimsWithoutValidation(tokenString string) (tokens.TokenClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Provider) GetProviderType() tokens.TokenProviderType {
	return tokens.TokenProviderJWE
}
