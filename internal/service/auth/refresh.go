package auth

import (
	"context"
	"errors"
	"fmt"
	"tgiintiz/internal/models"
)

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	claims, err := s.parse(refreshToken, refreshTokenType)
	if err != nil {
		return nil, err
	}

	userID, err := claims.userID()
	if err != nil {
		return nil, err
	}

	active, err := s.sessionRepo.Exists(ctx, claims.ID)
	if err != nil {
		return nil, fmt.Errorf("exists session: %w", err)
	}

	if !active {
		return nil, fmt.Errorf("%w: session revoked", models.ErrInvalidToken)
	}

	if err := s.sessionRepo.Delete(ctx, claims.ID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrInvalidToken
		}
		return nil, fmt.Errorf("delete session: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrInvalidToken
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	tokenPair, err := s.newTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}
