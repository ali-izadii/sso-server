package tokens

//
//import (
//	"context"
//	"fmt"
//	"sso-server/internal/domain/models"
//	"sso-server/internal/domain/repositories"
//	"sync"
//	"time"
//
//	"github.com/google/uuid"
//)
//
//type tokenManager struct {
//	providers       map[TokenProviderType]TokenProvider
//	defaultProvider TokenProviderType
//	tokenRepo       repositories.TokenRepository
//	eventListener   TokenEventListener
//	mu              sync.RWMutex
//}
//
//func NewTokenManager(tokenRepo repositories.TokenRepository) TokenManager {
//	return &tokenManager{
//		providers:       make(map[TokenProviderType]TokenProvider),
//		defaultProvider: TokenProviderJWT, // Default to JWT
//		tokenRepo:       tokenRepo,
//	}
//}
//
//func (tm *tokenManager) RegisterProvider(providerType TokenProviderType, provider TokenProvider) error {
//	tm.mu.Lock()
//	defer tm.mu.Unlock()
//
//	if provider == nil {
//		return fmt.Errorf("provider cannot be nil")
//	}
//
//	tm.providers[providerType] = provider
//	return nil
//}
//
//func (tm *tokenManager) GetProvider(providerType TokenProviderType) (TokenProvider, error) {
//	tm.mu.RLock()
//	defer tm.mu.RUnlock()
//
//	provider, exists := tm.providers[providerType]
//	if !exists {
//		return nil, fmt.Errorf("provider type %s not registered", providerType)
//	}
//
//	return provider, nil
//}
//
//func (tm *tokenManager) SetDefaultProvider(providerType TokenProviderType) error {
//	tm.mu.Lock()
//	defer tm.mu.Unlock()
//
//	if _, exists := tm.providers[providerType]; !exists {
//		return fmt.Errorf("provider type %s not registered", providerType)
//	}
//
//	tm.defaultProvider = providerType
//	return nil
//}
//
//func (tm *tokenManager) GetDefaultProvider() TokenProvider {
//	tm.mu.RLock()
//	defer tm.mu.RUnlock()
//
//	return tm.providers[tm.defaultProvider]
//}
//
//func (tm *tokenManager) CreateTokenPair(ctx context.Context, req models.CreateTokenRequest) (*models.TokenPair, error) {
//	return tm.CreateTokenPairWithProvider(ctx, tm.defaultProvider, req)
//}
//
//func (tm *tokenManager) CreateTokenPairWithProvider(ctx context.Context, providerType TokenProviderType, req models.CreateTokenRequest) (*models.TokenPair, error) {
//	provider, err := tm.GetProvider(providerType)
//	if err != nil {
//		return nil, err
//	}
//
//	// Generate access token
//	accessTokenString, accessClaims, err := provider.GenerateAccessToken(req)
//	if err != nil {
//		return nil, fmt.Errorf("failed to generate access token: %w", err)
//	}
//
//	// Convert to standard claims
//	standardAccessClaims := tm.convertToStandardClaims(accessClaims)
//
//	// Store access token
//	accessToken := &models.AccessToken{
//		ID:            standardAccessClaims.GetTokenID(),
//		Token:         accessTokenString,
//		UserID:        req.UserID,
//		ApplicationID: req.ApplicationID,
//		Scopes:        models.ScopesAsString(req.Scopes),
//		ExpiresAt:     standardAccessClaims.GetExpiresAt(),
//		Revoked:       false,
//		CreatedAt:     standardAccessClaims.GetIssuedAt(),
//	}
//
//	if err := tm.tokenRepo.CreateAccessToken(ctx, accessToken); err != nil {
//		return nil, fmt.Errorf("failed to store access token: %w", err)
//	}
//
//	// Generate refresh token
//	refreshTokenString, refreshClaims, err := provider.GenerateRefreshToken(req, accessToken.ID)
//	if err != nil {
//		// Cleanup access token
//		_ = tm.tokenRepo.RevokeAccessTokenByID(ctx, accessToken.ID)
//		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
//	}
//
//	// Convert to standard claims
//	standardRefreshClaims := tm.convertToStandardClaims(refreshClaims)
//
//	// Store refresh token
//	refreshToken := &models.RefreshToken{
//		ID:            standardRefreshClaims.GetTokenID(),
//		Token:         refreshTokenString,
//		UserID:        req.UserID,
//		ApplicationID: req.ApplicationID,
//		AccessTokenID: &accessToken.ID,
//		ExpiresAt:     standardRefreshClaims.GetExpiresAt(),
//		Revoked:       false,
//		CreatedAt:     standardRefreshClaims.GetIssuedAt(),
//	}
//
//	if err := tm.tokenRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
//		// Cleanup access token
//		_ = tm.tokenRepo.RevokeAccessTokenByID(ctx, accessToken.ID)
//		return nil, fmt.Errorf("failed to store refresh token: %w", err)
//	}
//
//	// Fire events
//	if tm.eventListener != nil {
//		tm.eventListener.OnTokenCreated(models.AccessTokenType, standardAccessClaims)
//		tm.eventListener.OnTokenCreated(models.RefreshTokenType, standardRefreshClaims)
//	}
//
//	// Build token pair
//	tokenPair := &models.TokenPair{
//		AccessToken: models.TokenResponse{
//			Token:     accessTokenString,
//			TokenType: "Bearer",
//			ExpiresAt: accessToken.ExpiresAt,
//			ExpiresIn: int64(accessToken.ExpiresAt.Sub(accessToken.CreatedAt).Seconds()),
//		},
//		RefreshToken: models.TokenResponse{
//			Token:     refreshTokenString,
//			TokenType: "Bearer",
//			ExpiresAt: refreshToken.ExpiresAt,
//			ExpiresIn: int64(refreshToken.ExpiresAt.Sub(refreshToken.CreatedAt).Seconds()),
//		},
//	}
//
//	return tokenPair, nil
//}
//
//func (tm *tokenManager) convertToStandardClaims(claims interface{}) TokenClaims {
//	// This would need to be implemented based on your specific claim types
//	// For now, return a wrapper that implements TokenClaims
//	return &standardClaims{data: claims}
//}
//
//type standardClaims struct {
//	data interface{}
//}
//
//func (sc *standardClaims) GetUserID() uuid.UUID {
//	// Implementation depends on your claim structure
//	return uuid.Nil
//}
//
//func (sc *standardClaims) GetApplicationID() uuid.UUID {
//	return uuid.Nil
//}
//
//func (sc *standardClaims) GetEmail() string {
//	return ""
//}
//
//func (sc *standardClaims) GetScopes() []string {
//	return []string{}
//}
//
//func (sc *standardClaims) GetTokenType() models.TokenType {
//	return models.AccessTokenType
//}
//
//func (sc *standardClaims) GetTokenID() uuid.UUID {
//	return uuid.Nil
//}
//
//func (sc *standardClaims) GetExpiresAt() time.Time {
//	return time.Time{}
//}
//
//func (sc *standardClaims) GetIssuedAt() time.Time {
//	return time.Time{}
//}
//
//func (sc *standardClaims) IsExpired() bool {
//	return time.Now().After(sc.GetExpiresAt())
//}
//
//func (sc *standardClaims) ToMap() map[string]interface{} {
//	return make(map[string]interface{})
//}
