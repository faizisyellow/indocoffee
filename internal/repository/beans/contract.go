package beans

import (
	"context"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type Beans interface {
	Insert(ctx context.Context, nw models.BeansModel) error
	GetAll(ctx context.Context) ([]models.BeansModel, error)
	GetById(ctx context.Context, id int) (models.BeansModel, error)
	Update(ctx context.Context, nw models.BeansModel) error
	Delete(ctx context.Context, id int) error
	DestroyMany(ctx context.Context) error
}

type Contract struct {
	NewBeans Beans
}

func (c *Contract) Test(t *testing.T) {

}
