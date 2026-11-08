package models

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type Pagination struct {
	Page     int
	PageSize int
}

func NewPagination(page, pageSize int) Pagination {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Pagination) Limit() int32 {
	return int32(p.PageSize)
}

func (p Pagination) Offset() int32 {
	return int32((p.Page - 1) * p.PageSize)
}
