package carts

import (
	"context"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type Carts interface {
	Insert(ctx context.Context, cart models.Cart) error
	GetById(ctx context.Context, cartid int) (models.Cart, error)
	IncrementQuantity(ctx context.Context, cartId int) error
	DecrementQuantity(ctx context.Context, cartId int) error
	Delete(ctx context.Context, cartId int) error
}

type Contract struct {
	NewCart Carts
}

func (c *Contract) Test(t *testing.T) {

}
