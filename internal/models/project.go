package models

import (
	"fmt"
	"strings"
	"time"
)

type Project struct {
	ID            int64
	Lang          string
	Slug          string
	Title         string
	Category      string
	Location      string
	CompletedYear *int16
	ShortDesc     string
	FullDesc      string
	ImageURL      string
	Images        []string
	IsActive      bool
	Client        string
	CreatedAt     time.Time
	Services      []*ServiceCard
}

type CreateProjectInput struct {
	Lang          string
	Slug          string
	Title         string
	Category      string
	Location      string
	CompletedYear *int16
	ShortDesc     string
	FullDesc      string
	ImageURL      string
	Images        []string
	IsActive      bool
	Client        string
}

type UpdateProjectInput struct {
	Slug          *string
	Title         *string
	Category      *string
	Location      *string
	CompletedYear *int16
	ShortDesc     *string
	FullDesc      *string
	ImageURL      *string
	Images        *[]string
	IsActive      *bool
	Client        *string
}

func (in *CreateProjectInput) Prepare() error {
	in.Lang = strings.TrimSpace(strings.ToLower(in.Lang))
	in.Slug = strings.TrimSpace(strings.ToLower(in.Slug))
	in.Title = strings.TrimSpace(in.Title)
	in.Category = strings.TrimSpace(in.Category)
	in.Location = strings.TrimSpace(in.Location)
	in.ShortDesc = strings.TrimSpace(in.ShortDesc)
	in.FullDesc = strings.TrimSpace(in.FullDesc)
	in.ImageURL = strings.TrimSpace(in.ImageURL)
	in.Client = strings.TrimSpace(in.Client)

	if in.Title == "" || in.Lang == "" || in.Slug == "" || in.ShortDesc == "" || in.Category == "" {
		return fmt.Errorf("%w: title, slug, lang, short_desc and category are required", ErrInvalidInput)
	}

	if !Lang(in.Lang).IsValid() {
		return fmt.Errorf("%w: invalid lang (allowed: ru, tj, en)", ErrInvalidInput)
	}

	if err := validateCompletedYear(in.CompletedYear); err != nil {
		return err
	}

	newImages, err := normalizeImages(in.Images)
	if err != nil {
		return err
	}
	in.Images = newImages

	return nil
}

func (in *UpdateProjectInput) ApplyTo(project *Project) error {
	if err := applyRequiredString(&project.Slug, in.Slug, "slug"); err != nil {
		return err
	}
	project.Slug = strings.ToLower(project.Slug)

	if err := applyRequiredString(&project.Title, in.Title, "title"); err != nil {
		return err
	}
	if err := applyRequiredString(&project.Category, in.Category, "category"); err != nil {
		return err
	}
	if err := applyRequiredString(&project.ShortDesc, in.ShortDesc, "short desc"); err != nil {
		return err
	}

	if in.Location != nil {
		project.Location = strings.TrimSpace(*in.Location)
	}
	if in.FullDesc != nil {
		project.FullDesc = strings.TrimSpace(*in.FullDesc)
	}
	if in.ImageURL != nil {
		project.ImageURL = strings.TrimSpace(*in.ImageURL)
	}
	if in.Client != nil {
		project.Client = strings.TrimSpace(*in.Client)
	}

	if in.CompletedYear != nil {
		if err := validateCompletedYear(in.CompletedYear); err != nil {
			return err
		}

		project.CompletedYear = in.CompletedYear
	}

	if in.Images != nil {
		newImages, err := normalizeImages(*in.Images)
		if err != nil {
			return err
		}

		project.Images = newImages

	}

	if in.IsActive != nil {
		project.IsActive = *in.IsActive
	}

	return nil
}

func normalizeImages(images []string) ([]string, error) {
	if images == nil {
		images = []string{}
	}

	if len(images) > 5 {
		return nil, fmt.Errorf("%w: maximum 5 images allowed", ErrInvalidInput)
	}

	if len(images) > 0 {
		cleaned := make([]string, 0, len(images))

		for _, img := range images {
			trimmed := strings.TrimSpace(img)
			if trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}

		images = cleaned
	}

	return images, nil
}

func validateCompletedYear(year *int16) error {
	if year == nil {
		return nil
	}

	if *year < 1950 || *year > 2100 {
		return fmt.Errorf("%w: invalid completed_year", ErrInvalidInput)
	}

	return nil
}

func applyRequiredString(dst, src *string, fieldName string) error {
	if src == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*src)
	if trimmed == "" {
		return fmt.Errorf("%w: %s cannot be empty", ErrInvalidInput, fieldName)
	}

	*dst = trimmed

	return nil
}
