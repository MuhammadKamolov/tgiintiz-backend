package models

import (
	"encoding/json"
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
