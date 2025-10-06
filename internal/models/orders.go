package models

import (
	"time"
)

type Order struct {
	Id                     string      `json:"id"`
	IdempotencyKey         string      `json:"idempotency_key"`
	CustomerId             int         `json:"customer_id"`
	CustomerName           string      `json:"customer_name"`
	CustomerEmail          string      `json:"customer_email"`
	Status                 string      `json:"status"`
	Items                  []OrderItem `json:"items"`
	TotalPrice             float64     `json:"total_price"`
	PhoneNumber            string      `json:"phone_number"`
	AlternativePhoneNumber *string     `json:"alternative_phone_number"`
	Street                 string      `json:"street"`
	City                   string      `json:"city"`
	CreatedAt              time.Time   `json:"created_at"`
}

type OrderItem struct {
	Id            int     `json:"id"`
	Image         string  `json:"image"`
	BeanName      string  `json:"bean_name"`
	FormName      string  `json:"form_name"`
	Roasted       string  `json:"roasted"`
	Price         float64 `json:"price"`
	OrderQuantity int     `json:"order_quantity"`
}
