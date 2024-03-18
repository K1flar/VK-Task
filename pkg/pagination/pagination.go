package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPageSize   = 10
	QueryPageName     = "page"
	QueryPageSizeName = "size"
)

type Pagination struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
}

func New(pageNumber, pageSize int) *Pagination {
	return &Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}

func (p *Pagination) ValidatePagination() {
	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}

	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
}

func (p *Pagination) GetLimit() int {
	return p.PageSize
}

func (p *Pagination) GetOffset() int {
	return (p.PageNumber - 1) * p.PageSize
}

func NewFromRequest(r *http.Request) *Pagination {
	pageNumber := parseInt(r.URL.Query().Get(QueryPageName), 1)
	pageSize := parseInt(r.URL.Query().Get(QueryPageSizeName), DefaultPageSize)
	return New(pageNumber, pageSize)
}

func parseInt(n string, def int) int {
	if n == "" {
		return def
	}
	num, err := strconv.Atoi(n)
	if err != nil {
		return def
	}
	return num
}
