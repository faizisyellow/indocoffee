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
