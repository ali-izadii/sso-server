package models

import "errors"

var (
	ErrTokenNotFound = errors.New("tokens not found")

	ErrTokenInvalid = errors.New("tokens is invalid")
	ErrTokenExpired = errors.New("tokens has expired")
	ErrTokenRevoked = errors.New("tokens has been revoked")

	ErrInvalidTokenType = errors.New("invalid tokens type")
	ErrWrongTokenType   = errors.New("wrong tokens type for this operation")

	ErrInvalidSignature = errors.New("invalid tokens signature")
	ErrInvalidClaims    = errors.New("invalid tokens claims")
	ErrMalformedToken   = errors.New("malformed tokens")
	ErrTokenNotYetValid = errors.New("tokens not yet valid")

	ErrTokenCreationFailed = errors.New("failed to create tokens")
	ErrTokenSigningFailed  = errors.New("failed to sign tokens")

	ErrRefreshTokenUsed     = errors.New("refresh tokens already used")
	ErrRefreshTokenMismatch = errors.New("refresh tokens does not match access tokens")

	ErrTooManyTokens = errors.New("too many active tokens for this user")

	ErrInsufficientScope = errors.New("insufficient scope for this operation")
	ErrInvalidScope      = errors.New("invalid scope")
)
