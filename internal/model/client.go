package model

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Secret       string    `json:"-"`
	RedirectURIs []string  `json:"redirect_uris"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateClientRequest struct {
	Name         string   `json:"name" validate:"required"`
	RedirectURIs []string `json:"redirect_uris" validate:"required,min=1"`
}

type ClientResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Secret       string    `json:"secret,omitempty"`
	RedirectURIs []string  `json:"redirect_uris"`
}
