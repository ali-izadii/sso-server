package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email address is not verified")

	ErrUnauthorized = errors.New("unauthorized")
)
