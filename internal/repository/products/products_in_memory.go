package products

import (
	"context"
	"errors"

	"github.com/faizisyellow/indocoffee/internal/models"
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
