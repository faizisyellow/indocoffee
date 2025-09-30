package products

import (
	"context"
	"errors"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type InMemoryProducts struct {
	Products []models.Product
}

func (p *InMemoryProducts) Insert(ctx context.Context, newProduct models.Product) error {

	if newProduct.Image == "" {
		return errors.New("errors image url is empty")
	}

	for _, product := range p.Products {
		if product.Roasted == newProduct.Roasted && product.FormId == newProduct.FormId && product.BeanId == newProduct.BeanId {
			return errors.New("product already exist")
		}
	}

	np := models.Product{
		Id:       newProduct.Id,
		Roasted:  newProduct.Roasted,
		Price:    newProduct.Price,
		Quantity: newProduct.Quantity,
		Image:    newProduct.Image,
		BeanId:   newProduct.BeanId,
		FormId:   newProduct.FormId,
	}

	p.Products = append(p.Products, np)

	return nil
}

func (p *InMemoryProducts) GetById(ctx context.Context, id int) (models.Product, error) {
	return models.Product{}, nil
}

func (p *InMemoryProducts) GetAll(ctx context.Context, qry repository.PaginatedProductsQuery) ([]models.Product, error) {
	return nil, nil
}

func (p *InMemoryProducts) Update(ctx context.Context, product models.Product) error {
	return nil
}

func (p *InMemoryProducts) DeleteMany(ctx context.Context) error {
	p.Products = nil
	return nil
}

func (p *InMemoryProducts) Delete(ctx context.Context, id int) error {

	return nil
}
