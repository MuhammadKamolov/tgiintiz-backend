package auth

import (
	"context"
	"errors"
	"fmt"
	"tgiintiz/internal/models"
)

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.parse(refreshToken, refreshTokenType)
	if err != nil {
		return err
	}

	if err := s.sessionRepo.Delete(ctx, claims.ID); err != nil && !errors.Is(err, models.ErrNotFound) {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
