package service

import (
	"context"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
	"tgiintiz/internal/service/auth"
	"tgiintiz/internal/service/user"
)

type Services struct {
	Auth AuthService
	User UserService
}

type AuthService interface {
	Register(ctx context.Context, in models.RegisterInput) (*models.User, error)
	Login(ctx context.Context, in models.LoginInput) (*models.TokenPair, error)
	ParseToken(token string) (*models.TokenPayload, error)
}

func NewServices(repos *repository.Repositories, tokenCfg auth.TokenConfig) *Services {
	authService := auth.NewService(repos.User, tokenCfg)
	userService := user.NewService(repos.User)

	return &Services{
		Auth: authService,
		User: userService,
	}
}

type UserService interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}
