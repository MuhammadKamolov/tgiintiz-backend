package auth

import (
	"context"
	"errors"
	"fmt"
	"tgiintiz/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(ctx context.Context, in models.LoginInput) (*models.TokenPair, error) {
	user, err := s.users.GetByEmail(ctx, normalizeEmail(in.Email))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	tokenPair, err := s.newTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("create tokens: %w", err)
	}

	return tokenPair, nil
}
