package carts

import (
	"context"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type Carts interface {
	Insert(ctx context.Context, cart models.Cart) error
}

type Contract struct {
	NewCart Carts
}

func (c *Contract) Test(t *testing.T) {

}
