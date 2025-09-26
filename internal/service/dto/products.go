package dto

type CreateProductMetadataRequest struct {
	Roasted  string  `json:"roasted" validate:"required,oneof=light medium dark"`
	Price    float64 `json:"price" validate:"required"`
	Quantity uint16  `json:"quantity" validate:"required"`
	Bean     uint    `json:"bean" validate:"required"`
	Form     uint    `json:"form" validate:"required"`
}

type UploadedFile struct {
	Name     string
	Size     int64
	Content  []byte
	MimeType string
}
