package jwt

import (
	"context"
	"fmt"
	"sso-server/pkg/tokens"
	"testing"
	"time"

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
		validateClaims func(t *testing.T, claims tokens.TokenClaims)
		validateToken  func(t *testing.T, tokenString string, config tokens.JWTConfig)
	}{
		{
			name: "successful token generation with default algorithm (HS256)",
			config: tokens.JWTConfig{
				AccessTokenExpiry: time.Hour,
				Issuer:            "test-issuer",
				Algorithm:         "",
				SecretKey:         "test-secret-key-at-least-256-bits",
			},
			request: tokens.CreateTokenRequest{
				UserID:        uuid.New(),
				ApplicationID: uuid.New(),
				Email:         "test@example.com",
				Scopes:        []string{"read", "write"},
			},
			expectedError: false,
			validateClaims: func(t *testing.T, claims tokens.TokenClaims) {
				assert.Equal(t, "test@example.com", claims.GetEmail())
				assert.Equal(t, []string{"read", "write"}, claims.GetScopes())
				assert.Equal(t, tokens.AccessTokenType, claims.GetTokenType())
				assert.False(t, claims.IsExpired())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			provider := Provider{config: test.config}
			ctx := context.Background()

			tokenString, claims, err := provider.GenerateAccessToken(ctx, test.request)

			if test.expectedError {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				assert.Nil(t, claims)
				assert.Contains(t, err.Error(), "failed to sign access token")
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, tokenString)
			require.NotNil(t, claims)

			fmt.Println(test.name)
		})
	}
}
