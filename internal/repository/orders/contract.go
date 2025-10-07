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
	UpdateOrdersStatus(ctx context.Context, orderId string, status OrderStatus) error
	UpdateOrdersStatusWithTx(ctx context.Context, tx *sql.Tx, orderId string, status OrderStatus) error
	GetOrderStatusById(ctx context.Context, orderId string) (string, error)
	GetOrderById(ctx context.Context, orderId string) (models.Order, error)
}

type Contract struct {
	NewOrder Orders
}

func (c *Contract) Test(t *testing.T) {

}
