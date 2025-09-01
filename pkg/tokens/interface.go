package tokens

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	AccessTokenType  TokenType = "access_token"
	RefreshTokenType TokenType = "refresh_token"
)

type CreateTokenRequest struct {
	UserID        uuid.UUID
	ApplicationID uuid.UUID
	Scopes        []string
	Email         string
}

type TokenClaims interface {
	GetUserID() uuid.UUID
	GetApplicationID() uuid.UUID
	GetEmail() string
	GetScopes() []string
	GetTokenType() TokenType
	GetTokenID() uuid.UUID
	GetExpiresAt() time.Time
	GetIssuedAt() time.Time
	IsExpired() bool
	ToMap() map[string]interface{}
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

type TokenValidationResult struct {
	Valid         bool
	Claims        TokenClaims
	Error         string
	ExpiresAt     time.Time
	TokenType     TokenType
	UserID        uuid.UUID
	ApplicationID uuid.UUID
	Scopes        []string
}

type TokenError struct {
	Code    string
	Message string
	Cause   error
}

const (
	ErrTokenExpired = "TOKEN_EXPIRED"
	ErrTokenInvalid = "TOKEN_INVALID"
	ErrTokenRevoked = "TOKEN_REVOKED"
)

type TokenProvider interface {
	// GenerateAccessToken Token Generation
	GenerateAccessToken(req CreateTokenRequest) (string, *TokenClaims, error)
	// GenerateRefreshToken  Refresh Token Generation
	GenerateRefreshToken(req CreateTokenRequest, accessTokenID uuid.UUID) (string, *TokenClaims, error)
	// ValidateToken Token Validation
	ValidateToken(tokenString string) (*TokenClaims, error)
	// ValidateAccessToken Access Token Validation
	ValidateAccessToken(tokenString string) (*TokenClaims, error)
	// ValidateRefreshToken  Refresh Token Validation
	ValidateRefreshToken(tokenString string) (*TokenClaims, error)
	// GetTokenInfo Token Information
	GetTokenInfo(tokenString string) (*TokenValidationResult, error)
	// ExtractClaimsWithoutValidation GetTokenInfo Token Information
	ExtractClaimsWithoutValidation(tokenString string) (*TokenClaims, error)
	// GetTokenExpiry Token Properties
	GetTokenExpiry(tokenType TokenType) time.Duration
	// GetProviderType Token Properties
	GetProviderType() TokenProviderType
}

type TokenProviderType string

const (
	TokenProviderJWT    TokenProviderType = "jwt"
	TokenProviderOpaque TokenProviderType = "opaque"
	TokenProviderPASETO TokenProviderType = "paseto"
	TokenProviderJWE    TokenProviderType = "jwe"
)

type TokenManager interface {
	// RegisterProvider Provider Management
	RegisterProvider(providerType TokenProviderType, provider TokenProvider) error
	// GetProvider RegisterProvider Provider Management
	GetProvider(providerType TokenProviderType) (TokenProvider, error)
	// SetDefaultProvider RegisterProvider Provider Management
	SetDefaultProvider(providerType TokenProviderType) error
	// GetDefaultProvider RegisterProvider Provider Management
	GetDefaultProvider() TokenProvider
	// CreateTokenPair Token Operations (uses default provider)
	CreateTokenPair(ctx context.Context, req CreateTokenRequest) (*TokenPair, error)
	// ValidateToken Token Operations (uses default provider)
	ValidateToken(ctx context.Context, tokenString string) (*TokenValidationResult, error)
	// RefreshTokens Token Operations (uses default provider)
	RefreshTokens(ctx context.Context, refreshTokenString string) (*TokenPair, error)
	// RevokeToken Token Operations (uses default provider)
	RevokeToken(ctx context.Context, tokenString string) error
	// CreateTokenPairWithProvider create a pair of tokens with a specific provider
	CreateTokenPairWithProvider(ctx context.Context, providerType TokenProviderType, req CreateTokenRequest) (*TokenPair, error)
	// ValidateTokenWithProvider validate a pair of tokens with a specific provider
	ValidateTokenWithProvider(ctx context.Context, providerType TokenProviderType, tokenString string) (*TokenValidationResult, error)
}

type TokenFactory interface {
	CreateJWTProvider(config JWTConfig) (TokenProvider, error)
	CreateOpaqueProvider(config OpaqueConfig) (TokenProvider, error)
	CreatePASETOProvider(config PASETOConfig) (TokenProvider, error)
}

type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
	Algorithm          string
}

type OpaqueConfig struct {
	TokenLength        int
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	DatabaseStore      OpaqueTokenStore
}

type PASETOConfig struct {
	SymmetricKey       []byte
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

type OpaqueTokenStore interface {
	StoreToken(ctx context.Context, token string, claims *TokenClaims, expiresAt time.Time) error
	GetTokenClaims(ctx context.Context, token string) (*TokenClaims, error)
	RevokeToken(ctx context.Context, token string) error
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
	CleanupExpiredTokens(ctx context.Context) (int64, error)
}

type TokenMetadata struct {
	ProviderType  TokenProviderType
	TokenType     TokenType
	CreatedAt     time.Time
	ExpiresAt     time.Time
	UserID        uuid.UUID
	ApplicationID uuid.UUID
	Revoked       bool
}

type TokenValidationOptions struct {
	SkipExpiryCheck    bool
	SkipSignatureCheck bool
	RequiredScopes     []string
	RequiredAudience   []string
	ClockSkew          time.Duration
}

type ExtendedTokenProvider interface {
	TokenProvider
	// ValidateTokenWithOptions Advanced validation
	ValidateTokenWithOptions(tokenString string, opts TokenValidationOptions) (*TokenClaims, error)
	// IntrospectToken Token introspection
	IntrospectToken(tokenString string) (*TokenMetadata, error)
	// ExtendTokenExpiry Token lifecycle
	ExtendTokenExpiry(tokenString string, additionalDuration time.Duration) error
	// RevokeTokenFamily Token revocation
	RevokeTokenFamily(tokenFamilyID uuid.UUID) error
	// ValidateMultipleTokens Batch operations
	ValidateMultipleTokens(tokenStrings []string) ([]*TokenValidationResult, error)
	// GenerateTokenBatch Batch operations
	GenerateTokenBatch(requests []CreateTokenRequest) ([]string, error)
}

type TokenEventListener interface {
	OnTokenCreated(tokenType TokenType, claims *TokenClaims)
	OnTokenValidated(tokenType TokenType, claims *TokenClaims)
	OnTokenExpired(tokenType TokenType, claims *TokenClaims)
	OnTokenRevoked(tokenType TokenType, claims *TokenClaims)
	OnTokenRefreshed(oldClaims, newClaims *TokenClaims)
}

type TokenAnalytics interface {
	GetTokenStats(ctx context.Context) (*TokenStats, error)
	GetTokenStatsByUser(ctx context.Context, userID uuid.UUID) (*UserTokenStats, error)
	GetTokenStatsByApplication(ctx context.Context, appID uuid.UUID) (*ApplicationTokenStats, error)
	GetTokenStatsByProvider(ctx context.Context, providerType TokenProviderType) (*ProviderTokenStats, error)
}

type TokenStats struct {
	TotalActive  int64
	TotalExpired int64
	TotalRevoked int64
	ByProvider   map[TokenProviderType]int64
	ByTokenType  map[TokenType]int64
}

type UserTokenStats struct {
	UserID       uuid.UUID
	ActiveTokens int64
	ByProvider   map[TokenProviderType]int64
	ByTokenType  map[TokenType]int64
	LastActivity time.Time
}

type ApplicationTokenStats struct {
	ApplicationID uuid.UUID
	ActiveTokens  int64
	ByProvider    map[TokenProviderType]int64
	ByTokenType   map[TokenType]int64
	LastActivity  time.Time
}

type ProviderTokenStats struct {
	ProviderType  TokenProviderType
	ActiveTokens  int64
	ByTokenType   map[TokenType]int64
	AverageExpiry time.Duration
}

type TokenRouter interface {
	DetectProviderType(tokenString string) (TokenProviderType, error)
	RouteValidation(tokenString string) (*TokenValidationResult, error)
	RouteIntrospection(tokenString string) (*TokenMetadata, error)
}
