package products

import (
	"context"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type Products interface {
	Insert(ctx context.Context, newProduct models.Product) error
}

type Contract struct {
	NewUsers func() (Products, func())
}

func (u Contract) Test(t *testing.T) {

}
