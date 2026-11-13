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

func (s *Service) UpdateStatus(ctx context.Context, id int64, status string) (*models.SurveyRequest, error) {
	if !models.SurveyStatus(status).IsValid() {
		return nil, fmt.Errorf("%w: unknown status %q", models.ErrInvalidInput, status)
	}

	survey, err := s.requestRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, fmt.Errorf("survey service: update status: %w", err)
	}

	return survey, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.requestRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("survey service: delete survey: %w", err)
	}

	return nil
}

func (s *Service) UploadAttachment(ctx context.Context, in models.UploadAttachmentInput) (*models.RequestAttachment, error) {
	survey, err := s.requestRepo.GetByID(ctx, in.RequestID)
	if err != nil {
		return nil, err
	}

	if in.UserRole != "admin" && survey.UserID != in.UserID {
		return nil, models.ErrForbidden
	}

	fileURL, err := s.storageRepo.Save(in.FileName, in.Content)
	if err != nil {
		return nil, fmt.Errorf("survey service: save file: %w", err)
	}

	newAttachment := models.RequestAttachment{
		RequestID: in.RequestID,
		FileName:  in.FileName,
		FileURL:   fileURL,
		FileSize:  in.FileSize,
	}

	att, err := s.attachmentRepo.Create(ctx, &newAttachment)
	if err != nil {
		_ = s.storageRepo.Delete(fileURL)
		return nil, fmt.Errorf("survey service: create attachment: %w", err)
	}

	return att, nil
}

func (s *Service) ListAttachments(ctx context.Context, requestID int64, claims models.TokenClaims) ([]*models.RequestAttachment, error) {
	survey, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("survey service: get survey: %w", err)
	}

	if claims.Role != "admin" && survey.UserID != claims.UserID {
		return nil, models.ErrForbidden
	}

	attachments, err := s.attachmentRepo.ListByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("survey service: list attachments: %w", err)
	}

	return attachments, nil
}

func (s *Service) DeleteAttachment(ctx context.Context, attachmentID int64, claims models.TokenClaims) error {
	att, err := s.attachmentRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("survey service: get attachment: %w", err)
	}

	survey, err := s.requestRepo.GetByID(ctx, att.RequestID)
	if err != nil {
		return fmt.Errorf("survey service: get survey: %w", err)
	}

	if claims.Role != "admin" && survey.UserID != claims.UserID {
		return models.ErrForbidden
	}

	if err := s.attachmentRepo.Delete(ctx, attachmentID); err != nil {
		return fmt.Errorf("survey service: delete attachment: %w", err)
	}

	if err := s.storageRepo.Delete(att.FileURL); err != nil {
		return fmt.Errorf("survey service: delete file url: %w", err)
	}

	return nil
}

func (s *Service) AddService(ctx context.Context, requestID, serviceID int64, claims models.TokenClaims) error {
	survey, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("survey service: get survey: %w", err)
	}

	if claims.Role != "admin" && survey.UserID != claims.UserID {
		return models.ErrForbidden
	}

	service, err := s.catalogRepo.GetByID(ctx, serviceID)
	if err != nil {
		return fmt.Errorf("servey service: get service by id: %w", err)
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

	if claims.Role != "admin" && survey.UserID != claims.UserID {
		return nil, models.ErrForbidden
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

	if claims.Role != "admin" && survey.UserID != claims.UserID {
		return models.ErrForbidden
	}

	if err := s.catalogRepo.RemoveByRequestID(ctx, requestID); err != nil {
		return fmt.Errorf("survey service: remove by req id: %w", err)
	}

	return nil
}
