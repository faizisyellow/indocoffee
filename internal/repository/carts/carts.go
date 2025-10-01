package carts

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

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
	query := `SELECT id,user_id,product_id,quantity,created_at FROM cart_items WHERE id = ?`

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
