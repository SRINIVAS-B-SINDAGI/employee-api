package repository

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
