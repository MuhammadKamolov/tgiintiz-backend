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
	survey     repository.SurveyRepository
	attachment repository.AttachmentRepository
	storage    storage.FileStorage
}

func New(survey repository.SurveyRepository, attachment repository.AttachmentRepository, storage storage.FileStorage) *Service {
	return &Service{
		survey:     survey,
		attachment: attachment,
		storage:    storage,
	}
}

func (s *Service) Create(ctx context.Context, userID int64, input models.SurveyRequestInput) (*models.SurveyRequest, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, fmt.Errorf("survey service: create: %w", models.ErrInvalidInput)
	}

	survey, err := s.survey.Create(ctx, userID, input)
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

func (s *Service) UploadAttachment(ctx context.Context, input models.UploadAttachmentInput) (*models.RequestAttachment, error) {
	survey, err := s.survey.GetByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	if input.UserRole != "admin" && survey.UserID != input.UserID {
		return nil, models.ErrForbidden
	}

	fileURL, err := s.storage.Save(input.FileName, input.Content)
	if err != nil {
		return nil, fmt.Errorf("survey service: save file: %w", err)
	}

	newAttachment := models.RequestAttachment{
		RequestID: input.RequestID,
		FileName:  input.FileName,
		FileURL:   fileURL,
		FileSize:  input.FileSize,
	}

	att, err := s.attachment.Create(ctx, &newAttachment)
	if err != nil {
		_ = s.storage.Delete(fileURL)
		return nil, fmt.Errorf("survey service: create attachment: %w", err)
	}

	return att, nil
}

func (s *Service) ListAttachments(ctx context.Context, requestID, userID int64, userRole string) ([]*models.RequestAttachment, error) {
	survey, err := s.survey.GetByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("survey service: get survey: %w", err)
	}

	if userRole != "admin" && survey.UserID != userID {
		return nil, models.ErrForbidden
	}

	attachments, err := s.attachment.ListByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("survey service: list attachments: %w", err)
	}

	return attachments, nil
}

func (s *Service) DeleteAttachment(ctx context.Context, attachmentID, userID int64, userRole string) error {
	att, err := s.attachment.GetByID(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("survey service: get attachment: %w", err)
	}

	survey, err := s.survey.GetByID(ctx, att.RequestID)
	if err != nil {
		return fmt.Errorf("survey service: get survey: %w", err)
	}

	if userRole != "admin" && survey.UserID != userID {
		return models.ErrForbidden
	}

	if err := s.attachment.Delete(ctx, attachmentID); err != nil {
		return fmt.Errorf("survey service: delete attachment: %w", err)
	}

	if err := s.storage.Delete(att.FileURL); err != nil {
		return fmt.Errorf("survey service: delete file url: %w", err)
	}

	return nil
}
