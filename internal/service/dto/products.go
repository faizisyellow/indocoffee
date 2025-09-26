package dto

type CreateProductMetadataRequest struct {
	Roasted  string  `json:"roasted" validate:"required,oneof=light medium dark"`
	Price    float64 `json:"price" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Bean     int     `json:"bean" validate:"required,min=1"`
	Form     int     `json:"form" validate:"required,min=1"`
}
