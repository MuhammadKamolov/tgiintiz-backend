package auth

import (
	"tgiintiz/internal/repository"
	"time"
)

type TokenConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Service struct {
	users      repository.UserRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(
	users repository.UserRepository,
	tokens TokenConfig,
) *Service {
	return &Service{
		users:      users,
		jwtSecret:  []byte(tokens.Secret),
		accessTTL:  tokens.AccessTTL,
		refreshTTL: tokens.RefreshTTL,
	}
}
