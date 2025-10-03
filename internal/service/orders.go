package service

import (
	"context"

	"github.com/faizisyellow/indocoffee/internal/service/dto"
)

type OrdersService struct {
}

func (o *OrdersService) Create(ctx context.Context, req dto.CreateOrderRequest, usrId int) error {

	// 1. Create payment
	// 2. Create order
	// 3. Delete cart
	return nil
}
