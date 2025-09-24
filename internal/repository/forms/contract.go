package forms

import (
	"context"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type Forms interface {
	Insert(ctx context.Context, nw models.FormsModel) error
	GetAll(ctx context.Context) ([]models.FormsModel, error)
	GetById(ctx context.Context, id int) (models.FormsModel, error)
	Update(ctx context.Context, nw models.FormsModel) error
	Delete(ctx context.Context, id int) error
	DestroyMany(ctx context.Context) error
}

type Contract struct {
	NewForms Forms
}

func (c *Contract) Test(t *testing.T) {

}
