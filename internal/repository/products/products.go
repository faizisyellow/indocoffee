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

func (p *ProductRepository) GetById(ctx context.Context, id int) (models.Product, error) {
	qry := `
	SELECT
       products.id,
       products.roasted,
       products.price,
       products.quantity,
       products.image,
       products.bean_id,
       products.form_id,
       beans.name  AS bean_name,
       forms.name  AS form_name
    FROM products
    JOIN beans ON beans.id = products.bean_id
    JOIN forms ON forms.id = products.form_id
    WHERE products.id = ?;
	`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	product := models.Product{}
	if err := p.Db.QueryRowContext(ctx, qry, id).Scan(
		&product.Id,
		&product.Roasted,
		&product.Price,
		&product.Quantity,
		&product.Image,
		&product.BeanId,
		&product.FormId,
		&product.BeansModel.Name,
		&product.FormsModel.Name,
	); err != nil {
		return models.Product{}, err
	}

	return product, nil
}
