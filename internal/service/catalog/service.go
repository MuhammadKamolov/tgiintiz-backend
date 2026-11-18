package catalog

import (
	"context"
	"fmt"
	"strings"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
)

type Service struct {
	catalogRepo repository.CatalogRepository
}

func New(catalogRepo repository.CatalogRepository) *Service {
	return &Service{catalogRepo: catalogRepo}
}

func (s *Service) Create(ctx context.Context, in models.CreateServiceInput) (*models.Service, error) {
	if err := in.Prepare(); err != nil {
		return nil, fmt.Errorf("catalog service: prepare create input: %w", err)
	}

	service, err := s.catalogRepo.Create(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("catalog service: create: %w", err)
	}

	return service, nil
}

func (s *Service) GetByID(ctx context.Context, serviceID int64) (*models.Service, error) {
	service, err := s.catalogRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("catalog service: get by serviceID: %w", err)
	}

	return service, nil
}

func (s *Service) GetBySlugAndLang(ctx context.Context, slug, lang string) (*models.Service, error) {
	slug = normalizeSlug(slug)
	lang = normalizeLang(lang)

	if slug == "" {
		return nil, fmt.Errorf("%w: slug is required", models.ErrInvalidInput)
	}

	if !models.Lang(lang).IsValid() {
		return nil, fmt.Errorf("%w: invalid lang (c)", models.ErrInvalidInput)
	}

	service, err := s.catalogRepo.GetBySlugAndLang(ctx, slug, lang)
	if err != nil {
		return nil, fmt.Errorf("catalog service: get by slug and lang: %w", err)
	}

	if !service.IsActive {
		return nil, models.ErrNotFound
	}

	return service, nil
}
func (s *Service) ListByLang(ctx context.Context, lang string) ([]*models.ServiceCard, error) {
	lang = normalizeLang(lang)

	if !models.Lang(lang).IsValid() {
		return nil, fmt.Errorf("%w: invalid lang (allowed: ru, tj, en)", models.ErrInvalidInput)
	}

	cards, err := s.catalogRepo.ListByLang(ctx, lang)
	if err != nil {
		return nil, fmt.Errorf("catalog service: list by lang: %w", err)
	}

	return cards, nil
}

func (s *Service) Update(ctx context.Context, serviceID int64, in models.UpdateServiceInput) (*models.Service, error) {
	service, err := s.catalogRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("catalog service: get by ID: %w", err)
	}

	if err := in.ApplyTo(service); err != nil {
		return nil, fmt.Errorf("catalog service: apply update: %w", err)
	}

	updatedService, err := s.catalogRepo.Update(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("catalog service: update: %w", err)
	}

	return updatedService, nil
}

func (s *Service) Delete(ctx context.Context, serviceID int64) error {
	if err := s.catalogRepo.Delete(ctx, serviceID); err != nil {
		return fmt.Errorf("catalog service: delete by serviceID: %w", err)
	}

	return nil
}

func normalizeLang(lang string) string {
	return strings.TrimSpace(strings.ToLower(lang))
}
func normalizeSlug(slug string) string {
	return strings.TrimSpace(strings.ToLower(slug))
}
