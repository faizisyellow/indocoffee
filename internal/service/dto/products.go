package dto

type CreateProductMetadataRequest struct {
	Roasted  string  `json:"roasted" validate:"required,oneof=light medium dark"`
	Price    float64 `json:"price" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Bean     int     `json:"bean" validate:"required,min=1"`
	Form     int     `json:"form" validate:"required,min=1"`
}

type GetProductResponse struct {
	Id       int     `json:"id"`
	Roasted  string  `json:"roasted"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Image    string  `json:"image"`
	BeanId   int     `json:"bean_id"`
	FormId   int     `json:"form_id"`
	Bean     struct {
		Name string `json:"name"`
	} `json:"bean"`

	Form struct {
		Name string `json:"name"`
	} `json:"form"`
}

type GetProductsResponse struct {
	GetProductResponse
}
