package repository

import (
	"strconv"
)

type PaginatedProductsQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
	Roast  string `json:"roast" validate:"light,medium,dark"`
	Form   int    `json:"form" validate:"gte=0"`
	Bean   int    `json:"bean" validate:"gte=0"`
}

type QueryProducts struct {
	Limit  string
	Offset string
	Sort   string
	Roast  string
	Form   string
	Bean   string
}

func (p PaginatedProductsQuery) Parse(r QueryProducts) (PaginatedProductsQuery, error) {
	limit := r.Limit
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return p, err
		}
		p.Limit = l
	}

	offset := r.Offset
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return p, err
		}
		p.Offset = o
	}

	sort := r.Sort
	if sort != "" {
		p.Sort = sort
	}

	roast := r.Roast
	if roast != "" {
		p.Roast = roast
	}

	form := r.Form
	if form != "" {
		f, err := strconv.Atoi(form)
		if err != nil {
			return p, err
		}
		p.Form = f
	}

	bean := r.Bean
	if bean != "" {
		b, err := strconv.Atoi(bean)
		if err != nil {
			return p, err
		}
		p.Form = b
	}

	return p, nil
}
