package products

import (
	"context"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type InMemoryProducts struct {
	Products []models.Product
}

func (p *InMemoryProducts) Insert(ctx context.Context, newProduct models.Product) error {

	return nil
}
