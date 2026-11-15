package catalog

import (
	"context"
	"encoding/json"
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
	if err := normalizeCreateInput(&in); err != nil {
		return nil, fmt.Errorf("catalog service: create: %w", err)
	}

	newService, err := s.catalogRepo.Create(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("catalog service: create: %w", err)
	}

	return newService, nil
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
		return nil, fmt.Errorf("slug are required: %w", models.ErrInvalidInput)
	}

	if !models.Lang(lang).IsValid() {
		return nil, fmt.Errorf("%w: invalid lang (allowed: ru, tj, en)", models.ErrInvalidInput)
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
		return nil, fmt.Errorf("catalog service: get by serviceID: %w", err)
	}

	if err := applyServiceUpdate(service, in); err != nil {
		return nil, fmt.Errorf("catalog service: apply service update: %w", err)
	}

	updated, err := s.catalogRepo.Update(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("catalog service: update: %w", err)
	}

	return updated, nil
}
func (s *Service) Delete(ctx context.Context, serviceID int64) error {
	if err := s.catalogRepo.Delete(ctx, serviceID); err != nil {
		return fmt.Errorf("catalog service: delete by serviceID: %w", err)
	}

	return nil
}

func normalizeCreateInput(in *models.CreateServiceInput) error {
	in.Lang = normalizeLang(in.Lang)
	in.Slug = normalizeSlug(in.Slug)
	in.Title = strings.TrimSpace(in.Title)
	in.ShortDesc = strings.TrimSpace(in.ShortDesc)
	in.FullDesc = strings.TrimSpace(in.FullDesc)
	in.ImageURL = strings.TrimSpace(in.ImageURL)

	var err error
	if in.Standards, err = normalizeJSONB(&in.Standards, "standards"); err != nil {
		return err
	}
	if in.Equipment, err = normalizeJSONB(&in.Equipment, "equipment"); err != nil {
		return err
	}

	if in.Title == "" || in.Lang == "" || in.Slug == "" || in.ShortDesc == "" {
		return fmt.Errorf("title, slug, lang and short_desc are required: %w", models.ErrInvalidInput)
	}

	if !models.Lang(in.Lang).IsValid() {
		return fmt.Errorf("invalid lang (allowed: ru, tj, en): %ws", models.ErrInvalidInput)
	}

	return nil
}

func normalizeLang(lang string) string {
	return strings.TrimSpace(strings.ToLower(lang))
}

func normalizeSlug(slug string) string {
	return strings.TrimSpace(strings.ToLower(slug))
}

func applyServiceUpdate(service *models.Service, in models.UpdateServiceInput) error {
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		if title == "" {
			return fmt.Errorf("title desc are required: %w", models.ErrInvalidInput)
		}
		service.Title = title
	}

	if in.ShortDesc != nil {
		shortDesc := strings.TrimSpace(*in.ShortDesc)
		if shortDesc == "" {
			return fmt.Errorf("short desc are required: %w", models.ErrInvalidInput)
		}
		service.ShortDesc = shortDesc
	}

	if in.FullDesc != nil {
		service.FullDesc = strings.TrimSpace(*in.FullDesc)
	}

	if in.ImageURL != nil {
		service.ImageURL = strings.TrimSpace(*in.ImageURL)
	}

	if in.Standards != nil {
		val, err := normalizeJSONB(in.Standards, "standards")
		if err != nil {
			return err
		}
		service.Standards = val
	}

	if in.Equipment != nil {
		val, err := normalizeJSONB(in.Equipment, "equipment")
		if err != nil {
			return err
		}
		service.Equipment = val
	}

	if in.IsActive != nil {
		service.IsActive = *in.IsActive
	}

	return nil
}

func normalizeJSONB(data *json.RawMessage, fieldName string) (json.RawMessage, error) {
	if data == nil || len(*data) == 0 {
		return json.RawMessage("[]"), nil
	}

	if !json.Valid(*data) {
		return nil, fmt.Errorf("%w: %s must be valid JSON", models.ErrInvalidInput, fieldName)
	}

	return *data, nil
}
