package auth

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"tgiintiz/internal/models"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLen   = 8
	maxPasswordBytes = 72
)

func (s *Service) Register(ctx context.Context, in models.RegisterInput) (*models.User, error) {
	in.Email = normalizeEmail(in.Email)
	if err := validateRegister(in); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, in.Email, string(hash))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateRegister(in models.RegisterInput) error {
	addr, err := mail.ParseAddress(in.Email)
	if err != nil || addr.Address != in.Email {
		return fmt.Errorf("%w: invalid email", models.ErrInvalidInput)
	}

	if utf8.RuneCountInString(in.Password) < minPasswordLen {
		return fmt.Errorf("%w: password must be at least %d characters", models.ErrInvalidInput, minPasswordLen)
	}

	if len(in.Password) > maxPasswordBytes {
		return fmt.Errorf("%w: password must be at most %d bytes", models.ErrInvalidInput, maxPasswordBytes)
	}

	return nil
}
