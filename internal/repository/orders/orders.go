package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type OrderStatus int

const (
	Confirm OrderStatus = iota
	Roasting
	Shipped
	Complete
	Cancelled
)

func (o OrderStatus) String() string {

	return []string{"confirm", "roasting", "shipped", "complete", "cancelled"}[o]
}

type OrdersRepository struct {
	Db *sql.DB
}

func (o *OrdersRepository) Create(ctx context.Context, tx *sql.Tx, newOrder models.Order) (string, error) {
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
			city,
			cart_ids,
			created_at
		) VALUES(?,?,?,?,?,?,CAST(? AS JSON),?,?,?,?,?,?)
	`

	itemsJSON, err := json.Marshal(newOrder.Items)
	if err != nil {
		return "", fmt.Errorf("failed to marshal items: %w", err)
	}

	cartIdsJSON, err := json.Marshal(newOrder.CartIds)
	if err != nil {
		return "", fmt.Errorf("failed to marshal items: %w", err)
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
		string(cartIdsJSON),
		time.Now().UTC(),
	)
	if err != nil {
		return "", err
	}

	return newOrder.Id, nil

}

func (o *OrdersRepository) GetIdempotencyKey(ctx context.Context, idemKey string) (string, error) {

	query := `SELECT idempotency_key FROM orders WHERE idempotency_key = ?`

	var idempotencyKey string

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	err := o.Db.QueryRowContext(ctx, query, idemKey).Scan(&idempotencyKey)
	return idempotencyKey, err
}

func (o *OrdersRepository) GetOrderStatusById(ctx context.Context, orderId string) (string, error) {
	query := `SELECT status FROM orders WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	var statusOrder string
	err := o.Db.QueryRowContext(ctx, query, orderId).Scan(&statusOrder)

	return statusOrder, err
}

func (o *OrdersRepository) UpdateOrdersStatus(ctx context.Context, orderId string, status OrderStatus) error {
	query := `UPDATE orders SET status  = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := o.Db.ExecContext(ctx, query, status.String(), orderId)
	return err
}

func (o *OrdersRepository) UpdateOrdersStatusWithTx(ctx context.Context, tx *sql.Tx, orderId string, status OrderStatus) error {
	query := `UPDATE orders SET status  = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, status.String(), orderId)
	return err
}

func (o *OrdersRepository) GetOrderById(ctx context.Context, orderId string) (models.Order, error) {
	query := `
		SELECT
			id,
			idempotency_key,
			customer_id,
			customer_name,
			customer_email,
			status,
			total_price,
			phone_number,
			alternative_phone_number,
			street,
			city,
			created_at,
			items,
			cart_ids
		FROM
		 orders where id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	var (
		order       models.Order
		itemsJSON   sql.NullString
		cartIdsJSON sql.NullString
	)

	err := o.Db.QueryRowContext(ctx, query, orderId).Scan(
		&order.Id,
		&order.IdempotencyKey,
		&order.CustomerId,
		&order.CustomerName,
		&order.CustomerEmail,
		&order.Status,
		&order.TotalPrice,
		&order.PhoneNumber,
		&order.AlternativePhoneNumber,
		&order.Street,
		&order.City,
		&order.CreatedAt,
		&itemsJSON,
		&cartIdsJSON,
	)

	if err != nil {
		return order, err
	}

	// Unmarshal only if not NULL
	if itemsJSON.Valid && itemsJSON.String != "" {
		err = json.Unmarshal([]byte(itemsJSON.String), &order.Items)
		if err != nil {
			return order, fmt.Errorf("failed to unmarshal items: %w", err)
		}
	}

	if cartIdsJSON.Valid && cartIdsJSON.String != "" {
		err = json.Unmarshal([]byte(cartIdsJSON.String), &order.CartIds)
		if err != nil {
			return order, fmt.Errorf("failed to unmarshal items: %w", err)
		}
	}

	return order, nil
}

func (o *OrdersRepository) GetOrders(ctx context.Context, qry repository.PaginatedOrdersQuery) ([]models.Order, error) {
	query := `
		SELECT
			id,
			idempotency_key,
			customer_id,
			customer_name,
			customer_email,
			status,
			total_price,
			phone_number,
			alternative_phone_number,
			street,
			city,
			created_at,
			items,
			cart_ids
		FROM (
			SELECT
				id,
				idempotency_key,
				customer_id,
				customer_name,
				customer_email,
				status,
				total_price,
				phone_number,
				alternative_phone_number,
				street,
				city,
				created_at,
				items,
				cart_ids
			FROM orders
			WHERE status LIKE concat("%",?,"%")
		) AS filtered_orders
		ORDER BY created_at ` + qry.Sort + `
		LIMIT ?
		OFFSET ?
	`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	rowsResult, err := o.Db.QueryContext(ctx, query, qry.Status, qry.Limit, qry.Offset)
	if err != nil {
		return nil, err
	}
	defer rowsResult.Close()

	var orders []models.Order

	for rowsResult.Next() {
		var (
			order       models.Order
			itemsJSON   sql.NullString
			cartIdsJSON sql.NullString
		)

		if err := rowsResult.Scan(
			&order.Id,
			&order.IdempotencyKey,
			&order.CustomerId,
			&order.CustomerName,
			&order.CustomerEmail,
			&order.Status,
			&order.TotalPrice,
			&order.PhoneNumber,
			&order.AlternativePhoneNumber,
			&order.Street,
			&order.City,
			&order.CreatedAt,
			&itemsJSON,
			&cartIdsJSON,
		); err != nil {
			return nil, err
		}

		// Unmarshal only if not NULL
		if itemsJSON.Valid && itemsJSON.String != "" {
			err = json.Unmarshal([]byte(itemsJSON.String), &order.Items)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal items: %w", err)
			}
		}

		if cartIdsJSON.Valid && cartIdsJSON.String != "" {
			err = json.Unmarshal([]byte(cartIdsJSON.String), &order.CartIds)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal items: %w", err)
			}
		}

		orders = append(orders, order)
	}

	return orders, rowsResult.Err()
}

func (o *OrdersRepository) GetTotalUsersOrders(ctx context.Context, usrId int) (int, error) {
	qry := `SELECT COUNT(*) FROM orders WHERE customer_id = ?"`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	var total int

	err := o.Db.QueryRowContext(ctx, qry, usrId).Scan(&total)
	return total, err
}
