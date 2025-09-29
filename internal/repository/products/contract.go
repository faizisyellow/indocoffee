package products

import (
	"context"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type Products interface {
	Insert(ctx context.Context, newProduct models.Product) error
	GetById(ctx context.Context, id int) (models.Product, error)
	GetAll(ctx context.Context, r repository.PaginatedProductsQuery) ([]models.Product, error)
}

type Contract struct {
	NewUsers func() (Products, func())
}

func (u Contract) Test(t *testing.T) {
	t.Run("create new product", func(t *testing.T) {

		t.Run("success create new product", func(t *testing.T) {
			var (
				ctx               = context.Background()
				product, teardown = u.NewUsers()
			)
			t.Cleanup(teardown)

			newProduct := models.Product{
				Id:       1,
				Roasted:  "light",
				Price:    15.7,
				Quantity: 100,
				Image:    "arabica.jpeg",
				BeanId:   1,
				FormId:   1,
			}

			err := product.Insert(ctx, newProduct)
			if err != nil {
				t.Errorf("expected to be success but got error: %v", err.Error())
				return
			}
		})

		t.Run("failed create new product because the product already exist", func(t *testing.T) {
			var (
				ctx               = context.Background()
				product, teardown = u.NewUsers()
			)
			t.Cleanup(teardown)

			newProduct := models.Product{
				Id:       1,
				Roasted:  "light",
				Price:    15.7,
				Quantity: 100,
				Image:    "arabica.jpeg",
				BeanId:   1,
				FormId:   1,
			}

			product.Insert(ctx, newProduct)

			// recreate again
			err := product.Insert(ctx, newProduct)
			if err == nil {
				t.Error("expected to be error but got success")
				return
			}
		})
	})
}
