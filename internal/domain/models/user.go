package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	FirstName     *string    `json:"first_name" db:"first_name"`
	LastName      *string    `json:"last_name" db:"last_name"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin     *time.Time `json:"last_login" db:"last_login"`
}

type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=8"`
	FirstName *string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name" binding:"omitempty,min=1,max=100"`
}

type UserResponse struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	FirstName     *string    `json:"first_name"`
	LastName      *string    `json:"last_name"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		LastLogin:     u.LastLogin,
	}
}

type CreateUserParams struct {
	Email        string
	PasswordHash string
	FirstName    *string
	LastName     *string
}

func (req *CreateUserRequest) Validate() error {
	if req.Email == "" {
		return ErrInvalidEmail
	}
	if len(req.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (u *User) GetFullName() string {
	var fullName string
	if u.FirstName != nil {
		fullName = *u.FirstName
	}
	if u.LastName != nil {
		if fullName != "" {
			fullName += " "
		}
		fullName += *u.LastName
	}
	return fullName
}
