package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	ID        int64
	Lang      string
	Slug      string
	Title     string
	ShortDesc string
	FullDesc  string
	ImageURL  string
	Standards json.RawMessage
	Equipment json.RawMessage
	IsActive  bool
	CreatedAt time.Time
}

type ServiceCard struct {
	ID        int64
	Lang      string
	Slug      string
	Title     string
	ShortDesc string
	ImageURL  string
	CreatedAt time.Time
}

type CreateServiceInput struct {
	Lang      string
	Slug      string
	Title     string
	ShortDesc string
	FullDesc  string
	ImageURL  string
	Standards json.RawMessage
	Equipment json.RawMessage
	IsActive  bool
}

type UpdateServiceInput struct {
	Title     *string
	ShortDesc *string
	FullDesc  *string
	ImageURL  *string
	Standards *json.RawMessage
	Equipment *json.RawMessage
	IsActive  *bool
}

type Lang string

const (
	LangRU      Lang = "ru"
	LangEN      Lang = "en"
	LangTJ      Lang = "tj"
	DefaultLang Lang = "ru"
)

func (l Lang) IsValid() bool {
	switch l {
	case LangRU, LangEN, LangTJ:
		return true
	default:
		return false
	}
}

func (in *CreateServiceInput) Prepare() error {
	in.Lang = strings.TrimSpace(strings.ToLower(in.Lang))
	in.Slug = strings.TrimSpace(strings.ToLower(in.Slug))
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
		return fmt.Errorf("%w: title, slug, lang and short_desc are required", ErrInvalidInput)
	}

	if !Lang(in.Lang).IsValid() {
		return fmt.Errorf("%w: invalid lang (allowed: ru, tj, en)", ErrInvalidInput)
	}

	return nil
}

func (in *UpdateServiceInput) ApplyTo(service *Service) error {
	if in.Title != nil {
		val := strings.TrimSpace(*in.Title)
		if val == "" {
			return fmt.Errorf("%w: title cannot be empty", ErrInvalidInput)
		}

		service.Title = val
	}

	if in.ShortDesc != nil {
		val := strings.TrimSpace(*in.ShortDesc)
		if val == "" {
			return fmt.Errorf("%w: short_desc cannot be empty", ErrInvalidInput)
		}

		service.ShortDesc = val
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
		return nil, fmt.Errorf("%w: %s must be valid JSON", ErrInvalidInput, fieldName)
	}

	return *data, nil
}
