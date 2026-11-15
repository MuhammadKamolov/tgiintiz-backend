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
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtSecret   []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func New(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	tokens TokenConfig,
) *Service {
	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   []byte(tokens.Secret),
		accessTTL:   tokens.AccessTTL,
		refreshTTL:  tokens.RefreshTTL,
	}
}
