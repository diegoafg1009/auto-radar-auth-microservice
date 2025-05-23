package repositories

import (
	"context"

	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/domain"
)

type User interface {
	Create(ctx context.Context, user *domain.User) (id string, err error)
	GetByID(ctx context.Context, id string) (user *domain.User, err error)
	GetByEmail(ctx context.Context, email string) (user *domain.User, err error)
}
