package v1

import (
	"tgiintiz/internal/models"
	"time"
)

type createProjectRequest struct {
	Lang          string   `json:"lang" binding:"required"`
	Slug          string   `json:"slug" binding:"required"`
	Title         string   `json:"title" binding:"required"`
	Category      string   `json:"category" binding:"required"`
	Location      string   `json:"location"`
	CompletedYear *int16   `json:"completed_year"`
	ShortDesc     string   `json:"short_desc" binding:"required"`
	FullDesc      string   `json:"full_desc"`
	ImageURL      string   `json:"image_url"`
	Images        []string `json:"images"`
	IsActive      *bool    `json:"is_active"`
	Client        string   `json:"client"`
}

func (req createProjectRequest) toInput() models.CreateProjectInput {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return models.CreateProjectInput{
		Lang:          req.Lang,
		Slug:          req.Slug,
		Title:         req.Title,
		Category:      req.Category,
		Location:      req.Location,
		CompletedYear: req.CompletedYear,
		ShortDesc:     req.ShortDesc,
		FullDesc:      req.FullDesc,
		ImageURL:      req.ImageURL,
		Images:        req.Images,
		IsActive:      isActive,
		Client:        req.Client,
	}
}

type updateProjectRequest struct {
	Slug          *string   `json:"slug"`
	Title         *string   `json:"title"`
	Category      *string   `json:"category"`
	Location      *string   `json:"location"`
	CompletedYear *int16    `json:"completed_year"`
	ShortDesc     *string   `json:"short_desc"`
	FullDesc      *string   `json:"full_desc"`
	ImageURL      *string   `json:"image_url"`
	Images        *[]string `json:"images"`
	IsActive      *bool     `json:"is_active"`
	Client        *string   `json:"client"`
}

func (req updateProjectRequest) toInput() models.UpdateProjectInput {
	return models.UpdateProjectInput{
		Slug:          req.Slug,
		Title:         req.Title,
		Category:      req.Category,
		Location:      req.Location,
		CompletedYear: req.CompletedYear,
		ShortDesc:     req.ShortDesc,
		FullDesc:      req.FullDesc,
		ImageURL:      req.ImageURL,
		Images:        req.Images,
		IsActive:      req.IsActive,
		Client:        req.Client,
	}
}

type addProjectServiceRequest struct {
	ServiceID int64 `json:"service_id" binding:"required,gt=0"`
}

type projectResponse struct {
	ID            int64                 `json:"id"`
	Lang          string                `json:"lang"`
	Slug          string                `json:"slug"`
	Title         string                `json:"title"`
	Category      string                `json:"category"`
	Location      string                `json:"location"`
	CompletedYear *int16                `json:"completed_year"`
	ShortDesc     string                `json:"short_desc"`
	FullDesc      string                `json:"full_desc"`
	ImageURL      string                `json:"image_url"`
	Images        []string              `json:"images"`
	IsActive      bool                  `json:"is_active"`
	Client        string                `json:"client"`
	CreatedAt     time.Time             `json:"created_at"`
	Services      []serviceCardResponse `json:"services"`
}

func newProjectResponse(p *models.Project) projectResponse {
	services := make([]serviceCardResponse, 0, len(p.Services))
	for _, s := range p.Services {
		services = append(services, newServiceCardResponse(s))
	}

	return projectResponse{
		ID:            p.ID,
		Lang:          p.Lang,
		Slug:          p.Slug,
		Title:         p.Title,
		Category:      p.Category,
		Location:      p.Location,
		CompletedYear: p.CompletedYear,
		ShortDesc:     p.ShortDesc,
		FullDesc:      p.FullDesc,
		ImageURL:      p.ImageURL,
		Images:        p.Images,
		IsActive:      p.IsActive,
		Client:        p.Client,
		CreatedAt:     p.CreatedAt,
		Services:      services,
	}
}
