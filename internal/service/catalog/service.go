package catalog

import (
	"context"
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
	return nil, nil
}
func (s *Service) AddToRequest(ctx context.Context, requestID, serviceID int64) error {
	return nil
}
func (s *Service) GetByID(ctx context.Context, serviceID int64) (*models.Service, error) {
	return nil, nil
}
func (s *Service) GetBySlugAndLang(ctx context.Context, slug, lang string) (*models.Service, error) {
	return nil, nil
}
func (s *Service) ListByLang(ctx context.Context, lang string) ([]*models.ServiceCard, error) {
	return nil, nil
}
func (s *Service) ListByRequestID(ctx context.Context, requestID int64) ([]*models.ServiceCard, error) {
	return nil, nil
}
func (s *Service) Update(ctx context.Context, serviceID int64, in models.UpdateServiceInput) (*models.Service, error) {
	return nil, nil
}
func (s *Service) RemoveByRequestID(ctx context.Context, requestID int64) error {
	return nil
}
func (s *Service) Delete(ctx context.Context, serviceID int64) error {
	return nil
}
