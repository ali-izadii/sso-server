package models

import "errors"

// Domain errors
var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters long")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserInactive      = errors.New("user account is inactive")
	ErrEmailNotVerified  = errors.New("email address is not verified")

	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenInvalid       = errors.New("invalid token")

	// Authorization errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// General errors
	ErrInternal         = errors.New("internal server error")
	ErrInvalidInput     = errors.New("invalid input")
	ErrResourceNotFound = errors.New("resource not found")
)
