package user

import (
	"context"
	"fmt"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
)

type Service struct {
	user repository.UserRepository
}

func NewService(user repository.UserRepository) *Service {
	return &Service{user: user}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.user.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user service: get by id: %w", err)
	}

	return user, nil
}
