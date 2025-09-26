package products

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type ProductRepository struct {
	Db *sql.DB
}

func (p *ProductRepository) Insert(ctx context.Context, newProduct models.Product) error {

	qry := `INSERT INTO products(roasted,price,quantity,image,bean_id,form_id) VALUE(?,?,?,?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := p.Db.ExecContext(
		ctx,
		qry,
		newProduct.Roasted,
		newProduct.Price,
		newProduct.Quantity,
		newProduct.Image,
		newProduct.BeanId,
		newProduct.FormId,
	)

	return err
}
