package project

import (
	"context"
	"errors"
	"fmt"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const uniqueViolation = "23505"
const keyViolation = "23503"

type repository struct {
	q *sqlc.Queries
}

func New(q *sqlc.Queries) *repository {
	return &repository{q: q}
}

func toModel(p sqlc.Project) *models.Project {
	images := p.Images
	if images == nil {
		images = []string{}
	}

	return &models.Project{
		ID:            p.ID,
		Slug:          p.Slug,
		Lang:          p.Lang,
		Title:         p.Title,
		Category:      p.Category,
		Location:      p.Location.String,
		CompletedYear: fromPgInt2(p.CompletedYear),
		ShortDesc:     p.ShortDesc,
		FullDesc:      p.FullDesc.String,
		ImageURL:      p.ImageUrl.String,
		Images:        images,
		IsActive:      p.IsActive,
		Client:        p.Client.String,
		CreatedAt:     p.CreatedAt.Time,
		Services:      []*models.ServiceCard{},
	}
}

func toCardModel(s sqlc.ListProjectServicesRow) *models.ServiceCard {
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

func (r *repository) AddService(ctx context.Context, projectID, serviceID int64) error {
	args := sqlc.AddServiceToProjectParams{
		ProjectID: projectID,
		ServiceID: serviceID,
	}

	if err := r.q.AddServiceToProject(ctx, args); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == keyViolation {
			return models.ErrNotFound
		}

		return fmt.Errorf("project repository: add service to project: %w", err)
	}

	return nil
}

func (r *repository) RemoveService(ctx context.Context, projectID, serviceID int64) error {
	args := sqlc.DeleteServiceFromProjectParams{
		ProjectID: projectID,
		ServiceID: serviceID,
	}

	if err := r.q.DeleteServiceFromProject(ctx, args); err != nil {
		return fmt.Errorf("project repository: remove service: %w", err)
	}

	return nil
}

func (r *repository) Create(ctx context.Context, in models.CreateProjectInput) (*models.Project, error) {
	args := sqlc.CreateProjectParams{
		Lang:          in.Lang,
		Slug:          in.Slug,
		Title:         in.Title,
		Category:      in.Category,
		Location:      toText(in.Location),
		CompletedYear: toPgInt2(in.CompletedYear),
		ShortDesc:     in.ShortDesc,
		FullDesc:      toText(in.FullDesc),
		ImageUrl:      toText(in.ImageURL),
		Images:        in.Images,
		IsActive:      in.IsActive,
		Client:        toText(in.Client),
	}

	row, err := r.q.CreateProject(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return nil, models.ErrSlugTaken
		}

		return nil, fmt.Errorf("project repository: create project: %w", err)
	}

	return toModel(row), nil
}
func (r *repository) Update(ctx context.Context, p *models.Project) (*models.Project, error) {
	args := sqlc.UpdateProjectParams{
		ID:            p.ID,
		Slug:          p.Slug,
		Title:         p.Title,
		Category:      p.Category,
		Location:      toText(p.Location),
		CompletedYear: toPgInt2(p.CompletedYear),
		ShortDesc:     p.ShortDesc,
		FullDesc:      toText(p.FullDesc),
		ImageUrl:      toText(p.ImageURL),
		Images:        p.Images,
		IsActive:      p.IsActive,
		Client:        toText(p.Client),
	}

	row, err := r.q.UpdateProject(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return nil, models.ErrSlugTaken
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("project repository: update project: %w", err)
	}

	return toModel(row), nil
}
func (r *repository) Delete(ctx context.Context, projectID int64) error {
	if err := r.q.DeleteProject(ctx, projectID); err != nil {
		return fmt.Errorf("project repository: delete project: %w", err)
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, projectID int64) (*models.Project, error) {
	row, err := r.q.GetProjectByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("project repository: get project by id: %w", err)
	}

	return toModel(row), nil
}
func (r *repository) GetBySlugAndLang(ctx context.Context, slug, lang string) (*models.Project, error) {
	args := sqlc.GetProjectBySlugAndLangParams{
		Slug: slug,
		Lang: lang,
	}

	row, err := r.q.GetProjectBySlugAndLang(ctx, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("project repository: get project by slug and lang: %w", err)
	}

	return toModel(row), nil
}
func (r *repository) ListActiveByLang(ctx context.Context, lang string, limit, offset int32) ([]*models.Project, error) {
	args := sqlc.ListActiveProjectsByLangParams{
		Lang:   lang,
		Limit:  limit,
		Offset: offset,
	}

	rows, err := r.q.ListActiveProjectsByLang(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("project repository: get active project lang: %w", err)
	}

	projects := make([]*models.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, toModel(row))
	}

	return projects, nil
}
func (r *repository) ListServicesByProjectID(ctx context.Context, projectID int64) ([]*models.ServiceCard, error) {
	rows, err := r.q.ListProjectServices(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project repository: list project services: %w", err)
	}

	cards := make([]*models.ServiceCard, 0, len(rows))
	for _, row := range rows {
		cards = append(cards, toCardModel(row))
	}

	return cards, nil
}
func (r *repository) ListByServiceID(ctx context.Context, serviceID int64, lang string) ([]*models.Project, error) {
	args := sqlc.ListProjectsByServiceParams{
		ServiceID: serviceID,
		Lang:      lang,
	}

	rows, err := r.q.ListProjectsByService(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("project repository: list projects by service: %w", err)
	}

	projects := make([]*models.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, toModel(row))
	}

	return projects, nil
}

func (r *repository) CountActiveByLang(ctx context.Context, lang string) (int64, error) {
	count, err := r.q.CountActiveProjectsByLang(ctx, lang)
	if err != nil {
		return 0, fmt.Errorf("project repository: count active projects: %w", err)
	}

	return count, nil
}

func toPgInt2(v *int16) pgtype.Int2 {
	if v == nil {
		return pgtype.Int2{}
	}

	return pgtype.Int2{Int16: *v, Valid: true}
}

func fromPgInt2(v pgtype.Int2) *int16 {
	if !v.Valid {
		return nil
	}

	val := v.Int16
	return &val
}

func toText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}
