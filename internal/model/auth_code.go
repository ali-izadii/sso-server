package model

import (
	"time"

	"github.com/google/uuid"
)

type AuthorizationCode struct {
	Code        string    `json:"code"`
	ClientID    uuid.UUID `json:"client_id"`
	UserID      uuid.UUID `json:"user_id"`
	RedirectURI string    `json:"redirect_uri"`
	Scope       string    `json:"scope"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}
