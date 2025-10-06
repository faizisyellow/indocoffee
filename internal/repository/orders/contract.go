package orders

import (
	"context"
	"database/sql"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type Orders interface {
	Create(ctx context.Context, tx *sql.Tx, nw models.Order) error
	GetIdempotencyKey(ctx context.Context, idemKey string) (string, error)
}

type Contract struct {
	NewOrder Orders
}

func (c *Contract) Test(t *testing.T) {

}
