package carts

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type CartStatus int

const (
	Open CartStatus = iota
	Ordered
)

func (c CartStatus) String() string {
	return []string{"open", "ordered"}[c]
}

type CartsRepository struct {
	Db *sql.DB
}

func (c *CartsRepository) Insert(ctx context.Context, cart models.Cart) error {
	query := `INSERT INTO cart_items(product_id,user_id) VALUES(?,?)`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := c.Db.ExecContext(ctx, query, cart.ProductId, cart.UserId)

	return err
}

func (c *CartsRepository) GetById(ctx context.Context, cartid int) (models.Cart, error) {
	query := `SELECT id,user_id,product_id,quantity,created_at FROM cart_items WHERE id = ? AND status = "open"`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	var cart models.Cart

	return cart, c.Db.QueryRowContext(ctx, query, cartid).Scan(
		&cart.Id,
		&cart.UserId,
		&cart.ProductId,
		&cart.Quantity,
		&cart.CreatedAt,
	)
}

func (c *CartsRepository) IncrementQuantity(ctx context.Context, cartid int) error {
	query := `UPDATE cart_items SET quantity = quantity + 1 WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := c.Db.ExecContext(ctx, query, cartid)

	return err
}

func (c *CartsRepository) DecrementQuantity(ctx context.Context, cartid int) error {
	query := `UPDATE cart_items SET quantity = quantity - 1 WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := c.Db.ExecContext(ctx, query, cartid)

	return err
}

func (c *CartsRepository) Delete(ctx context.Context, cartid int) error {

	qry := `DELETE FROM cart_items WHERE id = ?`
	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := c.Db.ExecContext(ctx, qry, cartid)
	return err
}

func (c *CartsRepository) DeleteWithTx(ctx context.Context, tx *sql.Tx, cartId int) error {

	qry := `DELETE FROM cart_items WHERE id = ?`
	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, qry, cartId)
	return err
}

func (c *CartsRepository) UpdateCartStatus(ctx context.Context, tx *sql.Tx, cartId int, status CartStatus) error {
	qry := `UPDATE cart_items SET status = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, qry, status.String(), cartId)
	return err
}

func (c *CartsRepository) GetCartStatus(ctx context.Context, cartId int) (string, error) {
	qry := `SELECT status FROM cart_items WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	var status string

	err := c.Db.QueryRowContext(ctx, qry, cartId).Scan(&status)

	return status, err
}

func (c *CartsRepository) GetTotalUsersCarts(ctx context.Context, usrId int) (int, error) {
	qry := `SELECT COUNT(*) FROM cart_items WHERE user_id = ?  AND status ="open"`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	var total int

	err := c.Db.QueryRowContext(ctx, qry, usrId).Scan(&total)
	return total, err
}
