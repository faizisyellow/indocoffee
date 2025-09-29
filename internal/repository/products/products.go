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

func (p *ProductRepository) GetAll(ctx context.Context, qry repository.PaginatedProductsQuery) ([]models.Product, error) {
	query := `
		SELECT
			p.id,
			p.roasted,
			p.price,
			p.quantity,
			p.image,
			p.bean_id,
			p.form_id,
			p.bean_name,
			p.form_name
		FROM (
			SELECT
				products.id,
				products.roasted,
				products.price,
				products.quantity,
				products.image,
				products.bean_id,
				products.form_id,
				beans.name AS bean_name,
				forms.name AS form_name
			FROM products
			JOIN beans ON beans.id = products.bean_id
			JOIN forms ON forms.id = products.form_id
			WHERE 1=1
`

	args := []any{}

	// add filters if present
	if qry.Roast != "" {
		query += " AND products.roasted LIKE CONCAT('%', ?, '%')"
		args = append(args, qry.Roast)
	}

	if qry.Form > 0 {
		query += " AND products.form_id = ?"
		args = append(args, qry.Form)
	}

	if qry.Bean > 0 {
		query += " AND products.bean_id = ?"
		args = append(args, qry.Bean)
	}

	// close subquery
	query += `
			ORDER BY beans.name ` + qry.Sort + `
		) AS p
	`

	// add pagination only if limit > 0
	if qry.Limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, qry.Limit, qry.Offset)
	}

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	rows, err := p.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
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
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *ProductRepository) DeleteMany(ctx context.Context) error {
	query := `DELETE FROM products`

	_, err := p.Db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
