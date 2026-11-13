package user

import (
	"context"
	"fmt"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
)

type Service struct {
	userRepo repository.UserRepository
}

func New(userRepo repository.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user service: get by id: %w", err)
	}

	return user, nil
}
