package auth

import (
	"context"
	"tgiintiz/internal/models"
)

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	return nil, nil
}
