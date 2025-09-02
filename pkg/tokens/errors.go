package tokens

import "errors"

var (
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenInvalid      = errors.New("token invalid")
	ErrTokenRevoked      = errors.New("token revoked")
	ErrTokenExpiredError = errors.New("token expired")
	ErrTokenInvalidError = errors.New("token invalid")
	ErrTokenRevokedError = errors.New("token revoked")
	ErrTokenNotFound     = errors.New("token not found")
	ErrInvalidSignature  = errors.New("invalid token signature")
	ErrInvalidClaims     = errors.New("invalid token claims")
	ErrInvalidTokenType  = errors.New("invalid token type")
	ErrProviderNotFound  = errors.New("provider not found")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrInsufficientScope = errors.New("insufficient scope")
)
