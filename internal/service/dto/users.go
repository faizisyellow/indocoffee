package dto

import "time"

type GetUsersProfileResponse struct {
	Id        int       `json:"int"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUsersCartResponse struct {
	Id         int              `json:"id"`
	Username   string           `json:"username"`
	Carts      []CartItemDetail `json:"carts"`
	TotalPrice float64          `json:"total_price"`
}

type CartItemDetail struct {
	Id       int            `json:"id"`
	Quantity int            `json:"quantity"`
	Product  CartProductDTO `json:"product"`
}

type CartProductDTO struct {
	Roasted string      `json:"roasted"`
	Image   string      `json:"image"`
	Stock   int         `json:"stock"` // (quantity)
	Price   float64     `json:"price"`
	Bean    CartBeanDTO `json:"bean"`
	Form    CartFormDTO `json:"form"`
}

type CartBeanDTO struct {
	Name string `json:"name"`
}

type CartFormDTO struct {
	Name string `json:"name"`
}
