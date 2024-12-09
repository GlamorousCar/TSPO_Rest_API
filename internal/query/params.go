package query

import (
	"net/http"
	"strconv"
)

type Params struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
	Filter   map[string]string
}

func NewParams(r *http.Request) *Params {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "title"
	}

	order := r.URL.Query().Get("order")
	if order != "desc" {
		order = "asc"
	}

	filter := make(map[string]string)
	if title := r.URL.Query().Get("title"); title != "" {
		filter["title"] = title
	}
	if author := r.URL.Query().Get("author"); author != "" {
		filter["author"] = author
	}

	return &Params{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
		Filter:   filter,
	}
}
