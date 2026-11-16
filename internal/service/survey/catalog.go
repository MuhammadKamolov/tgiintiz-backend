package survey

import (
	"context"
	"fmt"
	"tgiintiz/internal/models"
)

func (s *Service) AddService(ctx context.Context, requestID, serviceID int64, claims models.TokenClaims) error {
	survey, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("survey service: get survey: %w", err)
	}

	if err := checkSurveyAccess(survey, claims); err != nil {
		return fmt.Errorf("survey: service: check access: %w", err)
	}

	service, err := s.catalogRepo.GetByID(ctx, serviceID)
	if err != nil {
		return fmt.Errorf("survey service: add service: %w", err)
	}

	if !service.IsActive {
		return models.ErrNotFound
	}

	if err := s.catalogRepo.AddToRequest(ctx, requestID, serviceID); err != nil {
		return fmt.Errorf("survey service: add service: %w", err)
	}

	return nil
}
func (s *Service) ListServices(ctx context.Context, requestID int64, claims models.TokenClaims) ([]*models.ServiceCard, error) {
	survey, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("survey service: get survey: %w", err)
	}

	if err := checkSurveyAccess(survey, claims); err != nil {
		return nil, fmt.Errorf("survey: service: check access: %w", err)
	}

	card, err := s.catalogRepo.ListByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("survey service: list services: %w", err)
	}

	return card, nil
}
func (s *Service) RemoveServices(ctx context.Context, requestID int64, claims models.TokenClaims) error {
	survey, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("survey service: get survey: %w", err)
	}

	if err := checkSurveyAccess(survey, claims); err != nil {
		return fmt.Errorf("survey: service: check access: %w", err)
	}

	if err := s.catalogRepo.RemoveByRequestID(ctx, requestID); err != nil {
		return fmt.Errorf("survey service: remove by req id: %w", err)
	}

	return nil
}
