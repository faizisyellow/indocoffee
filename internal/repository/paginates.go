package repository

import (
	"strconv"
)

type PaginatedProductsQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0,lte=100"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
	Roast  string `json:"roast"`
	Form   int    `json:"form"`
	Bean   int    `json:"bean"`
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

type PaginatedOrdersQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0,lte=100"`
	Status string `json:"status" validate:"omitempty,oneof=confirm roasting shipped complete cancelled"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

type QueryOrders struct {
	Limit  string
	Offset string
	Sort   string
	Status string
}

func (p PaginatedOrdersQuery) Parse(r QueryOrders) (PaginatedOrdersQuery, error) {
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

	if r.Status != "" {
		p.Status = r.Status
	}

	return p, nil
}
