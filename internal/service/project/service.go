package project

import (
	"context"
	"fmt"
	"strings"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
)

type Service struct {
	projectRepo repository.ProjectRepository
}

func New(projectRepo repository.ProjectRepository) *Service {
	return &Service{projectRepo: projectRepo}
}

func (s *Service) AddService(ctx context.Context, projectID, serviceID int64) error {
	if err := s.projectRepo.AddService(ctx, projectID, serviceID); err != nil {
		return fmt.Errorf("project service: add service: %w", err)
	}

	return nil
}

func (s *Service) RemoveService(ctx context.Context, projectID, serviceID int64) error {
	if err := s.projectRepo.RemoveService(ctx, projectID, serviceID); err != nil {
		return fmt.Errorf("project service: remove service: %w", err)
	}

	return nil
}

func (s *Service) Create(ctx context.Context, in models.CreateProjectInput) (*models.Project, error) {
	if err := in.Prepare(); err != nil {
		return nil, fmt.Errorf("project service: prepare create input: %w", err)
	}

	project, err := s.projectRepo.Create(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("project service: create: %w", err)
	}

	return project, nil
}

func (s *Service) Update(ctx context.Context, projectID int64, in models.UpdateProjectInput) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project service: get by ID: %w", err)
	}

	if err := in.ApplyTo(project); err != nil {
		return nil, fmt.Errorf("project service: apply update: %w", err)
	}

	updatedProject, err := s.projectRepo.Update(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("project service: update: %w", err)
	}

	services, err := s.projectRepo.ListServicesByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project services: list services by project id: %w", err)
	}
	updatedProject.Services = services

	return updatedProject, nil
}

func (s *Service) Delete(ctx context.Context, projectID int64) error {
	if err := s.projectRepo.Delete(ctx, projectID); err != nil {
		return fmt.Errorf("project service: delete: %w", err)
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, projectID int64) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project service: get by ID: %w", err)
	}

	services, err := s.projectRepo.ListServicesByProjectID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("project service: list service by project id: %w", err)
	}
	project.Services = services

	return project, nil
}

func (s *Service) GetBySlugAndLang(ctx context.Context, slug, lang string) (*models.Project, error) {
	lang = normalizeLang(lang)
	if !models.Lang(lang).IsValid() {
		return nil, fmt.Errorf("%w: invalid lang (allowed: ru, tj, en)", models.ErrInvalidInput)
	}
	slug = normalizeSlug(slug)
	if slug == "" {
		return nil, fmt.Errorf("%w: slug is required", models.ErrInvalidInput)
	}

	project, err := s.projectRepo.GetBySlugAndLang(ctx, slug, lang)
	if err != nil {
		return nil, fmt.Errorf("project service: get by slug and lang: %w", err)
	}

	services, err := s.projectRepo.ListServicesByProjectID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("project service: list service by project id: %w", err)
	}
	project.Services = services

	return project, nil
}

func (s *Service) ListActiveByLang(ctx context.Context, lang string, p models.Pagination) ([]*models.Project, error) {
	lang = normalizeLang(lang)
	if !models.Lang(lang).IsValid() {
		return nil, fmt.Errorf("%w: invalid lang (allowed: ru, tj, en)", models.ErrInvalidInput)
	}

	projects, err := s.projectRepo.ListActiveByLang(ctx, lang, p.Limit(), p.Offset())
	if err != nil {
		return nil, fmt.Errorf("project service: list active by lang: %w", err)
	}

	return projects, nil
}

func (s *Service) ListServicesByProjectID(ctx context.Context, projectID int64) ([]*models.ServiceCard, error) {
	services, err := s.projectRepo.ListServicesByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project service: list services by project id: %w", err)
	}

	return services, nil
}

func (s *Service) ListByServiceID(ctx context.Context, serviceID int64, lang string) ([]*models.Project, error) {
	lang = normalizeLang(lang)
	if !models.Lang(lang).IsValid() {
		return nil, fmt.Errorf("%w: invalid lang (allowed: ru, tj, en)", models.ErrInvalidInput)
	}

	projects, err := s.projectRepo.ListByServiceID(ctx, serviceID, lang)
	if err != nil {
		return nil, fmt.Errorf("project service: list by services id: %w", err)
	}

	return projects, nil
}

func normalizeLang(lang string) string {
	return strings.TrimSpace(strings.ToLower(lang))
}
func normalizeSlug(slug string) string {
	return strings.TrimSpace(strings.ToLower(slug))
}

// генерация ключа токена через ssl rander
