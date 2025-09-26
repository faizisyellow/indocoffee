package products_test

import (
	"testing"

	"github.com/faizisyellow/indocoffee/internal/repository/products"
)

func TestProducts(t *testing.T) {
	products.Contract{func() (products.Products, func()) {
		return &products.InMemoryProducts{}, func() {}
	}}.Test(t)
}
