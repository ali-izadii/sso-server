package models

import "errors"

var (
	ErrInternal         = errors.New("internal server error")
	ErrInvalidInput     = errors.New("invalid input")
	ErrResourceNotFound = errors.New("resource not found")
	ErrForbidden        = errors.New("forbidden")
)
