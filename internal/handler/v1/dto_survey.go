package v1

import (
	"tgiintiz/internal/models"
	"time"
)

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

type surveyResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func newSurveyResponse(s *models.SurveyRequest) surveyResponse {
	return surveyResponse{
		ID:          s.ID,
		Title:       s.Title,
		Location:    s.Location,
		Description: s.Description,
		Status:      s.Status,
		CreatedAt:   s.CreatedAt,
	}
}

type updateSurveyStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_progress approved rejected"`
}
