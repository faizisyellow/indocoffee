package dto

type CreateOrderRequest struct {
	IndempotencyKey string            `json:"indempotency_key" validate:"required"`
	Name            string            `json:"name" validate:"required,min=1,max=32"`
	PhoneNumber     int               `json:"phone_number" validate:"required,min=1,max=18"`
	CartId          int               `json:"cart_id" validate:"required,min=1"`
	City            string            `json:"city" validate:"required,min=1,max=16"`
	Street          string            `json:"street" validate:"required,min=1,max=16"`
	Email           string            `json:"email" validate:"required,min=1,max=16"`
	Items           []CreateOrderItem `json:"items" validate:"required,min=1,max=16"`
}

type CreateOrderItem struct {
	Quantity  int `json:"quantity" validate:"required,min=1,max=16"`
	ProductId int `json:"product_id" validate:"required,min=1"`
}
