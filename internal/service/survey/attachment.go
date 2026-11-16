package survey

import (
	"context"
	"fmt"
	"tgiintiz/internal/models"
)

func (s *Service) UploadAttachment(ctx context.Context, in models.UploadAttachmentInput) (*models.RequestAttachment, error) {
	survey, err := s.requestRepo.GetByID(ctx, in.RequestID)
	if err != nil {
		return nil, err
	}

	if err := checkSurveyAccess(survey, in.Claims); err != nil {
		return nil, fmt.Errorf("survey service: check access: %w", err)
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

	if err := checkSurveyAccess(survey, claims); err != nil {
		return nil, fmt.Errorf("survey service: check access: %w", err)
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
	if err := checkSurveyAccess(survey, claims); err != nil {
		return fmt.Errorf("survey service: check access: %w", err)
	}
	if err := s.attachmentRepo.Delete(ctx, attachmentID); err != nil {
		return fmt.Errorf("survey service: delete attachment: %w", err)
	}

	if err := s.storageRepo.Delete(att.FileURL); err != nil {
		return fmt.Errorf("survey service: delete file url: %w", err)
	}

	return nil
}
