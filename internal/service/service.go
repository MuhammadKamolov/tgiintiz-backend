package service

import (
	"context"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
	"tgiintiz/internal/service/auth"
)

type Services struct {
	Auth AuthService
}

type AuthService interface {
	Register(ctx context.Context, in models.RegisterInput) (*models.User, error)
	Login(ctx context.Context, in models.LoginInput) (*models.TokenPair, error)
}

func NewServices(repos *repository.Repositories, tokenCfg auth.TokenConfig) *Services {
	authService := auth.NewService(repos.User, tokenCfg)
	return &Services{
		Auth: authService,
	}
}
