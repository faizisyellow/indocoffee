package dto

import (
	"time"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type CreateOrderRequest struct {
	CartIds                []int   `json:"cart_ids" validate:"required,min=1,max=16"`
	CustomerName           string  `json:"customer_name" validate:"required,min=1,max=32"`
	CustomerEmail          string  `json:"customer_email" validate:"required,email,min=6,max=32"`
	PhoneNumber            string  `json:"phone_number" validate:"required,min=10,max=15"`
	AlternativePhoneNumber *string `json:"alternative_phone_number" validate:"omitempty,min=10,max=15"`
	City                   string  `json:"city" validate:"required,min=5,max=32"`
	Street                 string  `json:"street" validate:"required,min=5,max=32"`
}

type GetOrderResponse struct {
	Id                     string             `json:"id"`
	CustomerName           string             `json:"customer_name"`
	CustomerEmail          string             `json:"customer_email"`
	Status                 string             `json:"status"`
	Items                  []models.OrderItem `json:"items"`
	TotalPrice             float64            `json:"total_price"`
	PhoneNumber            string             `json:"phone_number"`
	AlternativePhoneNumber *string            `json:"alternative_phone_number"`
	Street                 string             `json:"street"`
	City                   string             `json:"city"`
	CreatedAt              time.Time          `json:"created_at"`
}
