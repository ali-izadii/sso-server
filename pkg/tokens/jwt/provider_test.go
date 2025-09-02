package jwt

import (
	"sso-server/pkg/tokens"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_GenerateAccessToken(t *testing.T) {
	tests := []struct {
		name           string
		config         tokens.JWTConfig
		request        tokens.CreateTokenRequest
		expectedError  bool
		errorContains  string
		validateClaims func(t *testing.T, claims tokens.TokenClaims, request tokens.CreateTokenRequest, config tokens.JWTConfig)
		validateToken  func(t *testing.T, tokenString string, config tokens.JWTConfig)
	}{
		{
			name: "successful token generation with default algorithm (HS256)",
			config: tokens.JWTConfig{
				AccessTokenExpiry: time.Hour,
				Issuer:            "test-issuer",
				Algorithm:         "",
				SecretKey:         []byte("test-secret-key-at-least-256-bits-long"),
			},
			request: tokens.CreateTokenRequest{
				UserID:        uuid.New(),
				ApplicationID: uuid.New(),
				Email:         "test@example.com",
				Scopes:        []string{"read", "write"},
			},
			expectedError: false,
			validateClaims: func(t *testing.T, claims tokens.TokenClaims, request tokens.CreateTokenRequest, config tokens.JWTConfig) {
				assert.Equal(t, request.Email, claims.GetEmail())
				assert.Equal(t, request.UserID, claims.GetUserID())
				assert.Equal(t, request.ApplicationID, claims.GetApplicationID())
				assert.Equal(t, request.Scopes, claims.GetScopes())
				assert.Equal(t, tokens.AccessTokenType, claims.GetTokenType())
				assert.False(t, claims.IsExpired())
				assert.NotEqual(t, uuid.Nil, claims.GetTokenID())

				// Validate timing
				issuedAt := claims.GetTokenIssuedAt()
				expiresAt := claims.GeTokenExpirationTime()
				assert.True(t, issuedAt.Before(expiresAt))
				assert.True(t, expiresAt.Sub(issuedAt) <= config.AccessTokenExpiry+time.Second)
				assert.True(t, expiresAt.Sub(issuedAt) >= config.AccessTokenExpiry-time.Second)
			},
			validateToken: func(t *testing.T, tokenString string, config tokens.JWTConfig) {
				// Parse token to ensure it's valid JWT
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return config.SecretKey, nil
				})
				require.NoError(t, err)
				assert.True(t, token.Valid)
				assert.Equal(t, "HS256", token.Header["alg"])
			},
		},
		{
			name: "successful token generation with HS384",
			config: tokens.JWTConfig{
				AccessTokenExpiry: 2 * time.Hour,
				Issuer:            "test-issuer-384",
				Algorithm:         "HS384",
				SecretKey:         []byte("test-secret-key-for-hs384-algorithm-needs-to-be-longer"),
			},
			request: tokens.CreateTokenRequest{
				UserID:        uuid.New(),
				ApplicationID: uuid.New(),
				Email:         "user384@example.com",
				Scopes:        []string{"admin", "read", "write"},
			},
			expectedError: false,
			validateToken: func(t *testing.T, tokenString string, config tokens.JWTConfig) {
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return config.SecretKey, nil
				})
				require.NoError(t, err)
				assert.True(t, token.Valid)
				assert.Equal(t, "HS384", token.Header["alg"])
			},
		},
		{
			name: "successful token generation with HS512",
			config: tokens.JWTConfig{
				AccessTokenExpiry: 30 * time.Minute,
				Issuer:            "test-issuer-512",
				Algorithm:         "HS512",
				SecretKey:         []byte("test-secret-key-for-hs512-algorithm-should-be-even-longer-than-others"),
			},
			request: tokens.CreateTokenRequest{
				UserID:        uuid.New(),
				ApplicationID: uuid.New(),
				Email:         "user512@example.com",
				Scopes:        []string{"limited"},
			},
			expectedError: false,
			validateToken: func(t *testing.T, tokenString string, config tokens.JWTConfig) {
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return config.SecretKey, nil
				})
				require.NoError(t, err)
				assert.True(t, token.Valid)
				assert.Equal(t, "HS512", token.Header["alg"])
			},
		},
		{
			name: "empty scopes should return empty slice",
			config: tokens.JWTConfig{
				AccessTokenExpiry: time.Hour,
				Issuer:            "test-issuer",
				Algorithm:         "HS256",
				SecretKey:         []byte("test-secret-key-at-least-256-bits-long"),
			},
			request: tokens.CreateTokenRequest{
				UserID:        uuid.New(),
				ApplicationID: uuid.New(),
				Email:         "noscopes@example.com",
				Scopes:        nil,
			},
			expectedError: false,
			validateClaims: func(t *testing.T, claims tokens.TokenClaims, request tokens.CreateTokenRequest, config tokens.JWTConfig) {
				assert.Empty(t, claims.GetScopes())
				assert.NotNil(t, claims.GetScopes()) // Should be empty slice, not nil
			},
		},
		{
			name: "unsupported algorithm should fail",
			config: tokens.JWTConfig{
				AccessTokenExpiry: time.Hour,
				Issuer:            "test-issuer",
				Algorithm:         "RS256", // Unsupported algorithm
				SecretKey:         []byte("test-secret-key"),
			},
			request: tokens.CreateTokenRequest{
				UserID:        uuid.New(),
				ApplicationID: uuid.New(),
				Email:         "test@example.com",
				Scopes:        []string{"read"},
			},
			expectedError: true,
			errorContains: "unsupported signing algorithm: RS256",
		},
		{
			name: "another unsupported algorithm should fail",
			config: tokens.JWTConfig{
				AccessTokenExpiry: time.Hour,
				Issuer:            "test-issuer",
				Algorithm:         "PS256", // Another unsupported algorithm
				SecretKey:         []byte("test-secret-key"),
			},
			request: tokens.CreateTokenRequest{
				UserID:        uuid.New(),
				ApplicationID: uuid.New(),
				Email:         "test@example.com",
				Scopes:        []string{"read"},
			},
			expectedError: true,
			errorContains: "unsupported signing algorithm: PS256",
		},
		{
			name: "verify claims mapping with custom issuer and long expiry",
			config: tokens.JWTConfig{
				AccessTokenExpiry: 24 * time.Hour,
				Issuer:            "custom-sso-service",
				Algorithm:         "HS256",
				SecretKey:         []byte("super-secret-key-for-production-use-much-longer"),
			},
			request: tokens.CreateTokenRequest{
				UserID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				ApplicationID: uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
				Email:         "admin@company.com",
				Scopes:        []string{"admin", "user:read", "user:write", "system:manage"},
			},
			expectedError: false,
			validateClaims: func(t *testing.T, claims tokens.TokenClaims, request tokens.CreateTokenRequest, config tokens.JWTConfig) {
				// Verify all request fields are properly mapped
				assert.Equal(t, request.UserID, claims.GetUserID())
				assert.Equal(t, request.ApplicationID, claims.GetApplicationID())
				assert.Equal(t, request.Email, claims.GetEmail())
				assert.Equal(t, request.Scopes, claims.GetScopes())

				// Convert to CustomJwtClaims to test JWT-specific fields
				customClaims, ok := claims.(*CustomJwtClaims)
				require.True(t, ok, "Claims should be of type *CustomJwtClaims")

				assert.Equal(t, config.Issuer, customClaims.Issuer)
				assert.Equal(t, request.UserID.String(), customClaims.Subject)
				assert.Equal(t, customClaims.TokenID.String(), customClaims.ID)
				assert.Contains(t, customClaims.Audience, request.ApplicationID.String())

				// Verify claims map
				claimsMap := claims.ToMap()
				assert.Equal(t, request.UserID.String(), claimsMap["user_id"])
				assert.Equal(t, request.ApplicationID.String(), claimsMap["application_id"])
				assert.Equal(t, request.Email, claimsMap["email"])
				assert.Equal(t, request.Scopes, claimsMap["scopes"])
				assert.Equal(t, string(tokens.AccessTokenType), claimsMap["token_type"])
				assert.Equal(t, config.Issuer, claimsMap["iss"])
				assert.Equal(t, request.UserID.String(), claimsMap["sub"])
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := Provider{config: test.config}

			// Call the actual method (note: no context parameter based on interface)
			tokenString, claims, err := provider.GenerateAccessToken(test.request)

			if test.expectedError {
				assert.Error(t, err)
				if test.errorContains != "" {
					assert.Contains(t, err.Error(), test.errorContains)
				}
				assert.Empty(t, tokenString)
				assert.Nil(t, claims)
				return
			}

			// Success case validations
			require.NoError(t, err)
			assert.NotEmpty(t, tokenString)
			require.NotNil(t, claims)

			// Run custom claim validations if provided
			if test.validateClaims != nil {
				test.validateClaims(t, claims, test.request, test.config)
			}

			// Run custom token validations if provided
			if test.validateToken != nil {
				test.validateToken(t, tokenString, test.config)
			}

			// Basic token format validation
			parts := strings.Split(tokenString, ".")
			assert.Len(t, parts, 3, "JWT should have 3 parts separated by dots")
		})
	}
}

// Additional helper test for edge cases
func TestProvider_GenerateAccessToken_EdgeCases(t *testing.T) {
	t.Run("zero UUIDs should be handled", func(t *testing.T) {
		provider := Provider{
			config: tokens.JWTConfig{
				AccessTokenExpiry: time.Hour,
				Issuer:            "test",
				SecretKey:         []byte("test-secret-key-long-enough-for-hs256"),
			},
		}

		request := tokens.CreateTokenRequest{
			UserID:        uuid.Nil,   // Zero UUID
			ApplicationID: uuid.Nil,   // Zero UUID
			Email:         "",         // Empty email
			Scopes:        []string{}, // Empty scopes
		}

		tokenString, claims, err := provider.GenerateAccessToken(request)

		require.NoError(t, err)
		assert.NotEmpty(t, tokenString)
		require.NotNil(t, claims)

		assert.Equal(t, uuid.Nil, claims.GetUserID())
		assert.Equal(t, uuid.Nil, claims.GetApplicationID())
		assert.Equal(t, "", claims.GetEmail())
		assert.Empty(t, claims.GetScopes())
	})

	t.Run("very long scopes and email should be handled", func(t *testing.T) {
		provider := Provider{
			config: tokens.JWTConfig{
				AccessTokenExpiry: time.Hour,
				Issuer:            "test",
				SecretKey:         []byte("test-secret-key-long-enough-for-hs256"),
			},
		}

		longEmail := "very-long-email-address-that-might-cause-issues@very-long-domain-name-example.com"
		longScopes := []string{
			"very:long:scope:name:that:might:cause:issues:in:token:generation",
			"another:very:long:scope:name:with:many:colons:and:parts",
			"third:extremely:long:scope:name:for:testing:purposes:only",
		}

		request := tokens.CreateTokenRequest{
			UserID:        uuid.New(),
			ApplicationID: uuid.New(),
			Email:         longEmail,
			Scopes:        longScopes,
		}

		tokenString, claims, err := provider.GenerateAccessToken(request)

		require.NoError(t, err)
		assert.NotEmpty(t, tokenString)
		require.NotNil(t, claims)

		assert.Equal(t, longEmail, claims.GetEmail())
		assert.Equal(t, longScopes, claims.GetScopes())
	})
}
