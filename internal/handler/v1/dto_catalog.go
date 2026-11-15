package v1

import (
	"encoding/json"
	"tgiintiz/internal/models"
	"time"
)

type createServiceRequest struct {
	Lang      string          `json:"lang" binding:"required"`
	Slug      string          `json:"slug" binding:"required"`
	Title     string          `json:"title" binding:"required"`
	ShortDesc string          `json:"short_desc" binding:"required"`
	FullDesc  string          `json:"full_desc"`
	ImageURL  string          `json:"image_url"`
	Standards json.RawMessage `json:"standards"`
	Equipment json.RawMessage `json:"equipment"`
	IsActive  *bool           `json:"is_active"`
}

func (req createServiceRequest) toInput() models.CreateServiceInput {
	isActive := true

	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return models.CreateServiceInput{
		Lang:      req.Lang,
		Slug:      req.Slug,
		Title:     req.Title,
		ShortDesc: req.ShortDesc,
		FullDesc:  req.FullDesc,
		ImageURL:  req.ImageURL,
		Standards: req.Standards,
		Equipment: req.Equipment,
		IsActive:  isActive,
	}
}

type updateServiceRequest struct {
	Title     *string          `json:"title"`
	ShortDesc *string          `json:"short_desc"`
	FullDesc  *string          `json:"full_desc"`
	ImageURL  *string          `json:"image_url"`
	Standards *json.RawMessage `json:"standards"`
	Equipment *json.RawMessage `json:"equipment"`
	IsActive  *bool            `json:"is_active"`
}

func (req updateServiceRequest) toInput() models.UpdateServiceInput {
	return models.UpdateServiceInput{
		Title:     req.Title,
		ShortDesc: req.ShortDesc,
		FullDesc:  req.FullDesc,
		ImageURL:  req.ImageURL,
		Standards: req.Standards,
		Equipment: req.Equipment,
		IsActive:  req.IsActive,
	}
}

type addSurveyServiceRequest struct {
	ServiceID int64 `json:"service_id" binding:"required,gt=0"`
}

type serviceResponse struct {
	ID        int64           `json:"id"`
	Lang      string          `json:"lang"`
	Slug      string          `json:"slug"`
	Title     string          `json:"title"`
	ShortDesc string          `json:"short_desc"`
	FullDesc  string          `json:"full_desc"`
	ImageURL  string          `json:"image_url"`
	Standards json.RawMessage `json:"standards"`
	Equipment json.RawMessage `json:"equipment"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
}

func newServiceResponse(s *models.Service) serviceResponse {
	return serviceResponse{
		ID:        s.ID,
		Lang:      s.Lang,
		Slug:      s.Slug,
		Title:     s.Title,
		ShortDesc: s.ShortDesc,
		FullDesc:  s.FullDesc,
		ImageURL:  s.ImageURL,
		Standards: s.Standards,
		Equipment: s.Equipment,
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt,
	}
}

type serviceCardResponse struct {
	ID        int64     `json:"id"`
	Lang      string    `json:"lang"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	ShortDesc string    `json:"short_desc"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

func newServiceCardResponse(s *models.ServiceCard) serviceCardResponse {
	return serviceCardResponse{
		ID:        s.ID,
		Lang:      s.Lang,
		Slug:      s.Slug,
		Title:     s.Title,
		ShortDesc: s.ShortDesc,
		ImageURL:  s.ImageURL,
		CreatedAt: s.CreatedAt,
	}
}
