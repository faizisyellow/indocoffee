package dto

type CreateOrderRequest struct {
	CartIds                []int   `json:"cart_ids" validate:"required,min=1,max=16"`
	CustomerName           string  `json:"customer_name" validate:"required,min=1,max=32"`
	CustomerEmail          string  `json:"customer_email" validate:"required,email,min=6,max=32"`
	PhoneNumber            string  `json:"phone_number" validate:"required,min=10,max=15"`
	AlternativePhoneNumber *string `json:"alternative_phone_number" validate:"omitempty,min=10,max=15"`
	City                   string  `json:"city" validate:"required,min=5,max=32"`
	Street                 string  `json:"street" validate:"required,min=5,max=32"`
}
