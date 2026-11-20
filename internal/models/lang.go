package models

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
