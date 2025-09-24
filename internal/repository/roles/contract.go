package roles

import (
	"context"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type Roles interface {
	Insert(ctx context.Context, nw models.RolesModel) error
	GetAll(ctx context.Context) ([]models.RolesModel, error)
	GetById(ctx context.Context, id int) (models.RolesModel, error)
	Update(ctx context.Context, nw models.RolesModel) error
	Delete(ctx context.Context, id int) error
	DestroyMany(ctx context.Context) error
}

type Contract struct {
	NewRoles Roles
}

func (c *Contract) Test(t *testing.T) {

}
