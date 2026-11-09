package attachment

import (
	"context"
	"errors"
	"fmt"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository/sqlc"

	"github.com/jackc/pgx/v5"
)

type repository struct {
	q *sqlc.Queries
}

func New(q *sqlc.Queries) *repository {
	return &repository{
		q: q,
	}
}

func toModel(s sqlc.RequestAttachment) *models.RequestAttachment {
	return &models.RequestAttachment{
		ID:        s.ID,
		RequestID: s.RequestID,
		FileName:  s.FileName,
		FileURL:   s.FileUrl,
		FileSize:  s.FileSize,
		CreatedAt: s.CreatedAt.Time,
	}
}

func (r *repository) Create(ctx context.Context, att *models.RequestAttachment) (*models.RequestAttachment, error) {
	row, err := r.q.CreateRequestAttachment(ctx, sqlc.CreateRequestAttachmentParams{
		RequestID: att.RequestID,
		FileName:  att.FileName,
		FileUrl:   att.FileURL,
		FileSize:  att.FileSize,
	})
	if err != nil {
		return nil, fmt.Errorf("attachment repository: create: %w", err)
	}

	return toModel(row), nil
}

func (r *repository) ListByRequestID(ctx context.Context, requestID int64) ([]*models.RequestAttachment, error) {
	rows, err := r.q.ListAttachmentsByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("attachment repository: list by request id: %w", err)
	}

	attachments := make([]*models.RequestAttachment, 0, len(rows))
	for _, row := range rows {
		attachments = append(attachments, toModel(row))
	}

	return attachments, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*models.RequestAttachment, error) {
	row, err := r.q.GetAttachmentByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("attachment repository: get by id: %w", err)
	}

	return toModel(row), nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	if err := r.q.DeleteRequestAttachment(ctx, id); err != nil {
		return fmt.Errorf("attachment repository: delete by id: %w", err)
	}

	return nil
}
