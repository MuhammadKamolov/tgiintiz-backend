package survey

import (
	"context"
	"fmt"
	"strings"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
	"tgiintiz/pkg/storage"
)

type Service struct {
	requestRepo    repository.SurveyRepository
	attachmentRepo repository.AttachmentRepository
	catalogRepo    repository.CatalogRepository
	storageRepo    storage.FileStorage
}

func New(requestRepo repository.SurveyRepository, attachmentRepo repository.AttachmentRepository, catalogRepo repository.CatalogRepository, storageRepo storage.FileStorage) *Service {
	return &Service{
		requestRepo:    requestRepo,
		attachmentRepo: attachmentRepo,
		catalogRepo:    catalogRepo,
		storageRepo:    storageRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID int64, in models.SurveyRequestInput) (*models.SurveyRequest, error) {
	if strings.TrimSpace(in.Title) == "" {
		return nil, fmt.Errorf("survey service: create: %w", models.ErrInvalidInput)
	}

	survey, err := s.requestRepo.Create(ctx, userID, in)
	if err != nil {
		return nil, fmt.Errorf("survey service: create: %w", err)
	}

	return survey, nil
}

func (s *Service) GetByID(ctx context.Context, requestID int64, userID int64) (*models.SurveyRequest, error) {

	survey, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("survey service: get by id: %w", err)
	}

	if survey.UserID != userID {
		return nil, fmt.Errorf("survey service: get by id: %w", models.ErrForbidden)
	}

	return survey, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID int64, p models.Pagination) ([]*models.SurveyRequest, error) {
	survey, err := s.requestRepo.ListByUserID(ctx, userID, p.Limit(), p.Offset())
	if err != nil {
		return nil, fmt.Errorf("survey service: list by user id: %w", err)
	}

	return survey, nil
}

func (s *Service) ListAll(ctx context.Context, p models.Pagination) ([]*models.SurveyRequest, error) {
	survey, err := s.requestRepo.ListAll(ctx, p.Limit(), p.Offset())
	if err != nil {
		return nil, fmt.Errorf("survey service: list all: %w", err)
	}

	return survey, nil
}

func (s *Service) ListByStatus(ctx context.Context, status string, p models.Pagination) ([]*models.SurveyRequest, error) {
	if !models.SurveyStatus(status).IsValid() {
		return nil, fmt.Errorf("%w: unknown status %q", models.ErrInvalidInput, status)
	}

	survey, err := s.requestRepo.ListByStatus(ctx, status, p.Limit(), p.Offset())
	if err != nil {
		return nil, fmt.Errorf("survey service: list by status: %w", err)
	}

	return survey, nil
}

func (s *Service) UpdateStatus(ctx context.Context, requestID int64, status string) (*models.SurveyRequest, error) {
	if !models.SurveyStatus(status).IsValid() {
		return nil, fmt.Errorf("%w: unknown status %q", models.ErrInvalidInput, status)
	}

	survey, err := s.requestRepo.UpdateStatus(ctx, requestID, status)
	if err != nil {
		return nil, fmt.Errorf("survey service: update status: %w", err)
	}

	return survey, nil
}

func (s *Service) Delete(ctx context.Context, requestID int64) error {
	if err := s.requestRepo.Delete(ctx, requestID); err != nil {
		return fmt.Errorf("survey service: delete survey: %w", err)
	}

	return nil
}

func checkSurveyAccess(survey *models.SurveyRequest, claims models.TokenClaims) error {
	if !claims.IsAdmin() && survey.UserID != claims.UserID {
		return models.ErrForbidden
	}
	return nil
}
