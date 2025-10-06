package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type OrdersRepository struct {
	Db *sql.DB
}

func (o *OrdersRepository) Create(ctx context.Context, tx *sql.Tx, newOrder models.Order) error {
	qry := `
		INSERT INTO orders(
			id,
			idempotency_key,
			customer_id,
			customer_email,
			customer_name,
			items,
			total_price,
			phone_number,
			alternative_phone_number,
			street,
			city
		) VALUES(?,?,?,?,?,?,CAST(? AS JSON),?,?,?,?)
	`

	itemsJSON, err := json.Marshal(newOrder.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err = tx.ExecContext(
		ctx,
		qry,
		newOrder.Id,
		newOrder.IdempotencyKey,
		newOrder.CustomerId,
		newOrder.CustomerEmail,
		newOrder.CustomerName,
		string(itemsJSON),
		newOrder.TotalPrice,
		newOrder.PhoneNumber,
		newOrder.AlternativePhoneNumber,
		newOrder.Street,
		newOrder.City,
	)

	return err

}

func (o *OrdersRepository) GetIdempotencyKey(ctx context.Context, idemKey string) (string, error) {

	query := `SELECT idempotency_key FROM orders WHERE idempotency_key = ?`

	var idempotencyKey string

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	err := o.Db.QueryRowContext(ctx, query, idemKey).Scan(&idempotencyKey)
	return idempotencyKey, err
}

func (o *OrdersRepository) UpdateOrdersStatus(ctx context.Context, tx *sql.Tx, orderId int, status OrderStatus) error {
	query := `UPDATE orders SET status  = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, status.String(), orderId)
	return err
}

type OrderStatus int

const (
	Confirm OrderStatus = iota
	Roasting
	Shipped
	Complete
	Cancelled
)

func (o OrderStatus) String() string {

	return []string{"confirm", "roasted", "shipped", "complete", "cancelled"}[o]
}
