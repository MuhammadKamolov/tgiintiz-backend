package survey

import (
	"context"
	"fmt"
	"strings"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
)

type Service struct {
	survey repository.SurveyRepository
}

func New(survey repository.SurveyRepository) *Service {
	return &Service{survey: survey}
}

func (s *Service) Create(ctx context.Context, userID int64, in models.SurveyRequestInput) (*models.SurveyRequest, error) {
	if strings.TrimSpace(in.Title) == "" {
		return nil, fmt.Errorf("survey service: create: %w", models.ErrInvalidInput)
	}

	survey, err := s.survey.Create(ctx, userID, in.Title, in.Location, in.Description)
	if err != nil {
		return nil, fmt.Errorf("survey service: create: %w", err)
	}

	return survey, nil
}

func (s *Service) GetByID(ctx context.Context, id int64, userID int64) (*models.SurveyRequest, error) {

	survey, err := s.survey.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("survey service: get by id: %w", err)
	}

	if survey.UserID != userID {
		return nil, fmt.Errorf("survey service: get by id: %w", models.ErrForbidden)
	}

	return survey, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID int64, p models.Pagination) ([]*models.SurveyRequest, error) {
	survey, err := s.survey.ListByUserID(ctx, userID, p.Limit(), p.Offset())
	if err != nil {
		return nil, fmt.Errorf("survey service: list by user id: %w", err)
	}

	return survey, nil
}

func (s *Service) ListAll(ctx context.Context, p models.Pagination) ([]*models.SurveyRequest, error) {
	survey, err := s.survey.ListAll(ctx, p.Limit(), p.Offset())
	if err != nil {
		return nil, fmt.Errorf("survey service: list all: %w", err)
	}

	return survey, nil
}

func (s *Service) ListByStatus(ctx context.Context, status string, p models.Pagination) ([]*models.SurveyRequest, error) {
	if !models.SurveyStatus(status).IsValid() {
		return nil, fmt.Errorf("%w: unknown status %q", models.ErrInvalidInput, status)
	}

	survey, err := s.survey.ListByStatus(ctx, status, p.Limit(), p.Offset())
	if err != nil {
		return nil, fmt.Errorf("survey service: list by status: %w", err)
	}

	return survey, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id int64, status string) (*models.SurveyRequest, error) {
	if !models.SurveyStatus(status).IsValid() {
		return nil, fmt.Errorf("%w: unknown status %q", models.ErrInvalidInput, status)
	}

	survey, err := s.survey.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, fmt.Errorf("survey service: update status: %w", err)
	}

	return survey, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.survey.Delete(ctx, id); err != nil {
		return fmt.Errorf("survey service: delete survey: %w", err)
	}

	return nil
}
