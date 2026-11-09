package v1

import (
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

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type updateSurveyStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_progress approved rejected"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(u *models.User) userResponse {
	return userResponse{
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
