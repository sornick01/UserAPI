package user

import (
	"context"
	"github.com/sornick01/UserAPI/models"
)

type UseCase interface {
	CreateUser(ctx context.Context, displayName, email string) (string, error)
	GetUser(ctx context.Context, userId string) (*models.User, error)
	DeleteUser(ctx context.Context, userId string) error
	SearchUser(ctx context.Context) (map[string]*models.User, error)
	UpdateUser(ctx context.Context, userId, newName string) error
}
