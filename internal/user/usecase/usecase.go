package usecase

import (
	"context"
	"github.com/sornick01/UserAPI/internal/user"
	"github.com/sornick01/UserAPI/models"
	"time"
)

type DefaultUC struct {
	repo user.Repo
}

func NewDefaultUC(repo user.Repo) *DefaultUC {
	return &DefaultUC{
		repo: repo,
	}
}

func (d *DefaultUC) CreateUser(ctx context.Context, displayName, email string) (string, error) {
	u := &models.User{
		CreatedAt:   time.Now(),
		DisplayName: displayName,
		Email:       email,
	}

	return d.repo.CreateUser(ctx, u)
}

func (d *DefaultUC) GetUser(ctx context.Context, userId string) (*models.User, error) {
	return d.repo.GetUser(ctx, userId)
}

func (d *DefaultUC) DeleteUser(ctx context.Context, userId string) error {
	return d.repo.DeleteUser(ctx, userId)
}

func (d *DefaultUC) SearchUser(ctx context.Context) (map[string]*models.User, error) {
	return d.repo.SearchUser(ctx)
}

func (d *DefaultUC) UpdateUser(ctx context.Context, userId, newName string) error {
	return d.repo.UpdateUser(ctx, userId, newName)
}
