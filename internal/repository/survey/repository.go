package survey

import (
	"context"
	"errors"
	"fmt"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type repository struct {
	q *sqlc.Queries
}

func New(q *sqlc.Queries) *repository {
	return &repository{
		q: q,
	}
}

func toModel(s sqlc.SurveyRequest) *models.SurveyRequest {
	return &models.SurveyRequest{
		ID:          s.ID,
		UserID:      s.UserID,
		Title:       s.Title,
		Location:    s.Location.String,
		Description: s.Description.String,
		Status:      s.Status,
		CreatedAt:   s.CreatedAt.Time,
		UpdatedAt:   s.UpdatedAt.Time,
	}
}

func (r *repository) Create(ctx context.Context, userID int64, input models.SurveyRequestInput) (*models.SurveyRequest, error) {
	args := sqlc.CreateSurveyRequestParams{
		UserID:      userID,
		Title:       input.Title,
		Location:    pgtype.Text{String: input.Location, Valid: input.Location != ""},
		Description: pgtype.Text{String: input.Description, Valid: input.Description != ""},
	}

	row, err := r.q.CreateSurveyRequest(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("survey repository: create: %w", err)
	}

	return toModel(row), nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*models.SurveyRequest, error) {
	row, err := r.q.GetSurveyRequestByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("survey repository: get by id: %w", err)
	}

	return toModel(row), nil
}

func (r *repository) ListByUserID(ctx context.Context, userID int64, limit, offset int32) ([]*models.SurveyRequest, error) {
	rows, err := r.q.ListSurveyRequestsByUserID(ctx, sqlc.ListSurveyRequestsByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("survey repository: list by user id: %w", err)
	}

	survey := make([]*models.SurveyRequest, 0, len(rows))
	for _, row := range rows {
		survey = append(survey, toModel(row))
	}

	return survey, nil
}

func (r *repository) ListAll(ctx context.Context, limit, ofset int32) ([]*models.SurveyRequest, error) {
	rows, err := r.q.ListSurveyRequests(ctx, sqlc.ListSurveyRequestsParams{
		Limit:  limit,
		Offset: ofset,
	})
	if err != nil {
		return nil, fmt.Errorf("survey repository: list all: %w", err)
	}

	survey := make([]*models.SurveyRequest, 0, len(rows))
	for _, row := range rows {
		survey = append(survey, toModel(row))
	}

	return survey, nil
}

func (r *repository) ListByStatus(ctx context.Context, status string, limit, ofset int32) ([]*models.SurveyRequest, error) {
	rows, err := r.q.ListSurveyRequestsByStatus(ctx, sqlc.ListSurveyRequestsByStatusParams{
		Status: status,
		Limit:  limit,
		Offset: ofset,
	})
	if err != nil {
		return nil, fmt.Errorf("survey repository: list by status: %w", err)
	}

	survey := make([]*models.SurveyRequest, 0, len(rows))
	for _, row := range rows {
		survey = append(survey, toModel(row))
	}

	return survey, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id int64, status string) (*models.SurveyRequest, error) {
	args := sqlc.UpdateSurveyRequestStatusParams{
		ID:     id,
		Status: status,
	}

	survey, err := r.q.UpdateSurveyRequestStatus(ctx, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("survey repository: update status: %w", err)
	}

	return toModel(survey), nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	if err := r.q.DeleteSurveyRequest(ctx, id); err != nil {
		return fmt.Errorf("survey repository: delete: %w", err)
	}

	return nil
}
