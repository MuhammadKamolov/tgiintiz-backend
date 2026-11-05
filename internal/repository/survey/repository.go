package survey

import (
	"context"
	"errors"
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

func (r *repository) Create(ctx context.Context, userID int64, title string, location string, description string) (*models.SurveyRequest, error) {
	args := sqlc.CreateSurveyRequestParams{
		UserID:      userID,
		Title:       title,
		Location:    pgtype.Text{String: location, Valid: location != ""},
		Description: pgtype.Text{String: description, Valid: description != ""},
	}

	row, err := r.q.CreateSurveyRequest(ctx, args)
	if err != nil {
		return nil, err
	}

	return toModel(row), nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*models.SurveyRequest, error) {
	row, err := r.q.GetSurveyRequestByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return toModel(row), nil
}

func (r *repository) ListByUserID(ctx context.Context, userID int64) ([]*models.SurveyRequest, error) {
	rows, err := r.q.ListSurveyRequestsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	survey := make([]*models.SurveyRequest, 0, len(rows))

	for _, row := range rows {
		survey = append(survey, toModel(row))
	}

	return survey, nil
}
