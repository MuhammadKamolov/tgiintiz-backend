package v1

import (
	"encoding/json"
	"tgiintiz/internal/models"
	"time"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (req registerRequest) toInput() models.RegisterInput {
	return models.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (req loginRequest) toInput() models.LoginInput {
	return models.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

type surveyRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Location    string `json:"location" binding:"omitempty,max=255"`
	Description string `json:"description"`
}

func (req surveyRequest) toInput() models.SurveyRequestInput {
	return models.SurveyRequestInput{
		Title:       req.Title,
		Location:    req.Location,
		Description: req.Description,
	}
}

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

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type updateSurveyStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_progress approved rejected"`
}

type registerResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func newRegisterResponse(u *models.User) registerResponse {
	return registerResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

type surveyResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func newSurveyResponse(u *models.SurveyRequest) surveyResponse {
	return surveyResponse{
		ID:          u.ID,
		Title:       u.Title,
		Location:    u.Location,
		Description: u.Description,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
	}
}

type tokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newTokenPairResponse(tokens *models.TokenPair) tokenPairResponse {
	return tokenPairResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
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
