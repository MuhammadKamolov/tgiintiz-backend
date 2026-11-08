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

type SurveyStatus string

const (
	StatusPending    SurveyStatus = "pending"
	StatusInProgress SurveyStatus = "in_progress"
	StatusApproved   SurveyStatus = "approved"
	StatusRejected   SurveyStatus = "rejected"
)

func (s SurveyStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusApproved, StatusRejected:
		return true
	default:
		return false
	}
}
