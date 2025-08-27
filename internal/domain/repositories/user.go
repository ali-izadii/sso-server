package repositories

import (
	"context"

	"sso-server/internal/domain/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, params models.CreateUserParams) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id uuid.UUID, user *models.User) (*models.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	VerifyEmail(ctx context.Context, id uuid.UUID) error
	DeactivateUser(ctx context.Context, id uuid.UUID) error
	ReactivateUser(ctx context.Context, id uuid.UUID) error
	EmailExists(ctx context.Context, email string) (bool, error)
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)
}
