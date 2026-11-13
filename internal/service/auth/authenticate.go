package auth

import (
	"context"
	"fmt"
	"tgiintiz/internal/models"
)

func (s *Service) Authenticate(ctx context.Context, token string) (*models.TokenClaims, error) {
	claims, err := s.parse(token, accessTokenType)
	if err != nil {
		return nil, fmt.Errorf("parse access token: %w", err)
	}

	userID, err := claims.userID()
	if err != nil {
		return nil, err
	}

	active, err := s.sessionRepo.Exists(ctx, claims.ID)
	if err != nil {
		return nil, fmt.Errorf("check session: %w", err)
	}

	if !active {
		return nil, fmt.Errorf("%w: session revoked", models.ErrInvalidToken)
	}

	return &models.TokenClaims{
		UserID: userID,
		Role:   claims.Role,
	}, nil
}
