package catalog

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
	return &repository{q: q}
}

func toModel(s sqlc.Service) *models.Service {
	return &models.Service{
		ID:        s.ID,
		Lang:      s.Lang,
		Slug:      s.Slug,
		Title:     s.Title,
		ShortDesc: s.ShortDesc,
		FullDesc:  s.FullDesc.String,
		ImageURL:  s.ImageUrl.String,
		Standards: s.Standards,
		Equipment: s.Equipment,
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt.Time,
	}
}

func toCardModel(s sqlc.ListServicesByLangRow) *models.ServiceCard {
	return &models.ServiceCard{
		ID:        s.ID,
		Lang:      s.Lang,
		Slug:      s.Slug,
		Title:     s.Title,
		ShortDesc: s.ShortDesc,
		ImageURL:  s.ImageUrl.String,
		CreatedAt: s.CreatedAt.Time,
	}
}

func (r *repository) Create(ctx context.Context, in models.CreateServiceInput) (*models.Service, error) {
	args := sqlc.CreateServiceParams{
		Lang:      in.Lang,
		Slug:      in.Slug,
		Title:     in.Title,
		ShortDesc: in.ShortDesc,
		FullDesc:  pgtype.Text{String: in.FullDesc, Valid: in.FullDesc != ""},
		ImageUrl:  pgtype.Text{String: in.ImageURL, Valid: in.ImageURL != ""},
		Standards: in.Standards,
		Equipment: in.Equipment,
		IsActive:  in.IsActive,
	}
	row, err := r.q.CreateService(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("catalog repository: create: %w", err)
	}

	return toModel(row), nil
}

func (r *repository) AddToRequest(ctx context.Context, requestID, serviceID int64) error {
	args := sqlc.AddServiceToRequestParams{
		RequestID: requestID,
		ServiceID: serviceID,
	}

	if err := r.q.AddServiceToRequest(ctx, args); err != nil {
		return fmt.Errorf("catalog repository: add: %w", err)
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, serviceID int64) (*models.Service, error) {
	row, err := r.q.GetServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("catalog repository: get by id: %w", err)
	}

	return toModel(row), nil
}
func (r *repository) GetBySlugAndLang(ctx context.Context, slug, lang string) (*models.Service, error) {
	args := sqlc.GetServiceBySlugAndLangParams{
		Slug: slug,
		Lang: lang,
	}

	row, err := r.q.GetServiceBySlugAndLang(ctx, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("catalog repository: get by slug and lang: %w", err)
	}

	return toModel(row), nil
}
func (r *repository) ListByLang(ctx context.Context, lang string) ([]*models.ServiceCard, error) {
	rows, err := r.q.ListServicesByLang(ctx, lang)
	if err != nil {
		return nil, fmt.Errorf("catalog repository: list by lang: %w", err)
	}

	cards := make([]*models.ServiceCard, 0, len(rows))
	for _, r := range rows {
		cards = append(cards, toCardModel(r))
	}

	return cards, nil
}
func (r *repository) ListByRequestID(ctx context.Context, requestID int64) ([]*models.ServiceCard, error) {
	rows, err := r.q.ListServicesByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("catalog repository: list by request id: %w", err)
	}

	cards := make([]*models.ServiceCard, 0, len(rows))

	for _, r := range rows {
		cards = append(cards, toCardModel(sqlc.ListServicesByLangRow(r)))
	}

	return cards, nil
}
func (r *repository) Update(ctx context.Context, service *models.Service) (*models.Service, error) {
	args := sqlc.UpdateServiceParams{
		ID:        service.ID,
		Title:     service.Title,
		ShortDesc: service.ShortDesc,
		FullDesc:  pgtype.Text{String: service.FullDesc, Valid: service.FullDesc != ""},
		ImageUrl:  pgtype.Text{String: service.ImageURL, Valid: service.ImageURL != ""},
		Standards: service.Standards,
		Equipment: service.Equipment,
		IsActive:  service.IsActive,
	}

	row, err := r.q.UpdateService(ctx, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("catalog repository: update: %w", err)
	}
	return toModel(row), nil
}
func (r *repository) RemoveByRequestID(ctx context.Context, requestID int64) error {
	if err := r.q.RemoveServicesByRequestID(ctx, requestID); err != nil {
		return fmt.Errorf("catalog repository: remove by request id: %w", err)
	}

	return nil
}
func (r *repository) Delete(ctx context.Context, serviceID int64) error {
	if err := r.q.DeleteService(ctx, serviceID); err != nil {
		return fmt.Errorf("catalog repository: delete: %w", err)
	}
	return nil
}
