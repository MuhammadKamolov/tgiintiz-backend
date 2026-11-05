package models

import "time"

type SurveyRequest struct {
	ID          int64
	UserID      int64
	Title       string
	Location    string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SurveyRequestInput struct {
	Title       string
	Location    string
	Description string
}
