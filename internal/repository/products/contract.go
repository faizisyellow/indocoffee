package products

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type Products interface {
	Insert(ctx context.Context, newProduct models.Product) error
	GetById(ctx context.Context, id int) (models.Product, error)
	GetAll(ctx context.Context, r repository.PaginatedProductsQuery) ([]models.Product, error)
	Update(ctx context.Context, product models.Product) error
	DecrementQuantity(ctx context.Context, tx *sql.Tx, productId, quantity int) error
	IncrementQuantity(ctx context.Context, tx *sql.Tx, productId, quantity int) error
	DeleteMany(ctx context.Context) error
	Delete(ctx context.Context, id int) error
}

type Contract struct {
	NewProducts func() (Products, func())
}

func (u Contract) Test(t *testing.T) {
	t.Run("create new product", func(t *testing.T) {

		t.Run("success create new product", func(t *testing.T) {
			var (
				ctx               = context.Background()
				product, teardown = u.NewProducts()
			)
			t.Cleanup(func() {
				product.DeleteMany(ctx)
				teardown()
			})

			newProduct := models.Product{
				Roasted:  "light",
				Price:    10.5,
				Quantity: 50,
				Image:    "light_arabica_grounded.jpeg",
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
				product, teardown = u.NewProducts()
			)

			t.Cleanup(func() {
				product.DeleteMany(ctx)
				teardown()
			})

			newProduct := models.Product{
				Roasted:  "light",
				Price:    15.7,
				Quantity: 50,
				Image:    "light_arabica_grounded.jpeg",
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

	t.Run("get all products", func(t *testing.T) {
		tests := []struct {
			name     string
			query    repository.PaginatedProductsQuery
			expected []models.Product
		}{
			{
				name: "get products with default query",
				query: repository.PaginatedProductsQuery{
					Sort: "asc",
				},
				expected: []models.Product{
					{Roasted: "light", Price: 10.5, Quantity: 50, Image: "light_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "medium", Price: 12.0, Quantity: 70, Image: "medium_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "dark", Price: 14.8, Quantity: 30, Image: "dark_arabica_whole.jpeg", BeanId: 1, FormId: 2,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "whole coffee beans"}},
					{Roasted: "light", Price: 15.2, Quantity: 120, Image: "light_robusta_grounded.jpeg", BeanId: 2, FormId: 1,
						BeansModel: models.BeansModel{Name: "robusta"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "medium", Price: 18.0, Quantity: 90, Image: "medium_robusta_grounded.jpeg", BeanId: 2, FormId: 1,
						BeansModel: models.BeansModel{Name: "robusta"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "dark", Price: 20.0, Quantity: 40, Image: "dark_robusta_whole.jpeg", BeanId: 2, FormId: 2,
						BeansModel: models.BeansModel{Name: "robusta"}, FormsModel: models.FormsModel{Name: "whole coffee beans"}},
					{Roasted: "light", Price: 25.5, Quantity: 200, Image: "light_arabica_whole_premium.jpeg", BeanId: 1, FormId: 2,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "whole coffee beans"}},
					{Roasted: "dark", Price: 30.0, Quantity: 10, Image: "dark_robusta_grounded_limited.jpeg", BeanId: 2, FormId: 1,
						BeansModel: models.BeansModel{Name: "robusta"}, FormsModel: models.FormsModel{Name: "grounded"}},
				},
			},
			{
				name: "get all products only medium roast",
				query: repository.PaginatedProductsQuery{
					Sort:  "asc",
					Roast: "medium",
				},
				expected: []models.Product{
					{Roasted: "medium", Price: 12.0, Quantity: 70, Image: "medium_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "medium", Price: 18.0, Quantity: 90, Image: "medium_robusta_grounded.jpeg", BeanId: 2, FormId: 1,
						BeansModel: models.BeansModel{Name: "robusta"}, FormsModel: models.FormsModel{Name: "grounded"}},
				},
			},
			{
				name: "get all products only light roast and grounded form",
				query: repository.PaginatedProductsQuery{
					Sort:  "asc",
					Roast: "light",
					Form:  1,
				},
				expected: []models.Product{
					{Roasted: "light", Price: 10.5, Quantity: 50, Image: "light_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "light", Price: 15.2, Quantity: 120, Image: "light_robusta_grounded.jpeg", BeanId: 2, FormId: 1,
						BeansModel: models.BeansModel{Name: "robusta"}, FormsModel: models.FormsModel{Name: "grounded"}},
				},
			},
			{
				name: "get first page of all products",
				query: repository.PaginatedProductsQuery{
					Sort:   "asc",
					Limit:  3,
					Offset: 0,
				},
				expected: []models.Product{
					{Roasted: "light", Price: 10.5, Quantity: 50, Image: "light_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "medium", Price: 12.0, Quantity: 70, Image: "medium_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "dark", Price: 14.8, Quantity: 30, Image: "dark_arabica_whole.jpeg", BeanId: 1, FormId: 2,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "whole coffee beans"}},
				},
			},
			{
				name: "get only products that medium roast with pagination",
				query: repository.PaginatedProductsQuery{
					Sort:   "asc",
					Limit:  5,
					Offset: 0,
					Roast:  "medium",
				},
				expected: []models.Product{
					{Roasted: "medium", Price: 12.0, Quantity: 70, Image: "medium_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "medium", Price: 18.0, Quantity: 90, Image: "medium_robusta_grounded.jpeg", BeanId: 2, FormId: 1,
						BeansModel: models.BeansModel{Name: "robusta"}, FormsModel: models.FormsModel{Name: "grounded"}},
				},
			},
			{
				name: "get products only robusta",
				query: repository.PaginatedProductsQuery{
					Sort: "asc",
					Bean: 1,
				},
				expected: []models.Product{
					{Roasted: "light", Price: 10.5, Quantity: 50, Image: "light_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "medium", Price: 12.0, Quantity: 70, Image: "medium_arabica_grounded.jpeg", BeanId: 1, FormId: 1,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "grounded"}},
					{Roasted: "dark", Price: 14.8, Quantity: 30, Image: "dark_arabica_whole.jpeg", BeanId: 1, FormId: 2,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "whole coffee beans"}},
					{Roasted: "light", Price: 25.5, Quantity: 200, Image: "light_arabica_whole_premium.jpeg", BeanId: 1, FormId: 2,
						BeansModel: models.BeansModel{Name: "arabica"}, FormsModel: models.FormsModel{Name: "whole coffee beans"}},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				product, teardown := u.NewProducts()
				t.Cleanup(func() {
					product.DeleteMany(ctx)
					teardown()
				})

				if err := createTestProduct(t, product); err != nil {
					t.Fatal(err)
				}

				products, err := product.GetAll(ctx, tc.query)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if diff := cmp.Diff(tc.expected, products, cmpopts.IgnoreFields(models.Product{}, "Id")); diff != "" {
					t.Errorf("mismatch (-expected +got):\n%s", diff)
				}
			})
		}
	})
}

func createTestProduct(t *testing.T, p Products) error {
	t.Helper()

	newProducts := struct {
		input []models.Product
	}{
		input: []models.Product{
			// Arabica, grounded
			{Roasted: "light", Price: 10.5, Quantity: 50, Image: "light_arabica_grounded.jpeg", BeanId: 1, FormId: 1},
			{Roasted: "medium", Price: 12.0, Quantity: 70, Image: "medium_arabica_grounded.jpeg", BeanId: 1, FormId: 1},

			// Arabica, whole beans
			{Roasted: "dark", Price: 14.8, Quantity: 30, Image: "dark_arabica_whole.jpeg", BeanId: 1, FormId: 2},

			// Robusta, grounded
			{Roasted: "light", Price: 15.2, Quantity: 120, Image: "light_robusta_grounded.jpeg", BeanId: 2, FormId: 1},
			{Roasted: "medium", Price: 18.0, Quantity: 90, Image: "medium_robusta_grounded.jpeg", BeanId: 2, FormId: 1},

			// Robusta, whole beans
			{Roasted: "dark", Price: 20.0, Quantity: 40, Image: "dark_robusta_whole.jpeg", BeanId: 2, FormId: 2},

			// Extra variations for testing price/quantity ranges
			{Roasted: "light", Price: 25.5, Quantity: 200, Image: "light_arabica_whole_premium.jpeg", BeanId: 1, FormId: 2},
			{Roasted: "dark", Price: 30.0, Quantity: 10, Image: "dark_robusta_grounded_limited.jpeg", BeanId: 2, FormId: 1},
		},
	}

	for _, nw := range newProducts.input {
		err := p.Insert(context.Background(), nw)
		if err != nil {
			return fmt.Errorf("failed inserting product %+v: %w", nw, err)
		}
	}

	return nil
}
