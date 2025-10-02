package models

import (
	"time"
)

type Cart struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	Product   Product   `json:"product"`
	UserId    int       `json:"user_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
