package models

import "errors"

var (
	ErrTokenNotFound = errors.New("token not found")

	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenRevoked = errors.New("token has been revoked")

	ErrInvalidTokenType = errors.New("invalid token type")
	ErrWrongTokenType   = errors.New("wrong token type for this operation")

	ErrInvalidSignature = errors.New("invalid token signature")
	ErrInvalidClaims    = errors.New("invalid token claims")
	ErrMalformedToken   = errors.New("malformed token")
	ErrTokenNotYetValid = errors.New("token not yet valid")

	ErrTokenCreationFailed = errors.New("failed to create token")
	ErrTokenSigningFailed  = errors.New("failed to sign token")

	ErrRefreshTokenUsed     = errors.New("refresh token already used")
	ErrRefreshTokenMismatch = errors.New("refresh token does not match access token")

	ErrTooManyTokens = errors.New("too many active tokens for this user")

	ErrInsufficientScope = errors.New("insufficient scope for this operation")
	ErrInvalidScope      = errors.New("invalid scope")
)
