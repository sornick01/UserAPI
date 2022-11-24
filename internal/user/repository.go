package user

import (
	"context"
	"github.com/sornick01/UserAPI/models"
)

type Repo interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	DeleteUser(ctx context.Context, userId string) error
	GetUser(ctx context.Context, userId string) (*models.User, error)
	SearchUser(ctx context.Context) (map[string]*models.User, error)
	UpdateUser(ctx context.Context, userId, newName string) error
}
