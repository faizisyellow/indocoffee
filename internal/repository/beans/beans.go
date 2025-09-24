package beans

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type BeansRepository struct {
	Db *sql.DB
}

func (Beans *BeansRepository) Insert(ctx context.Context, nw models.BeansModel) error {

	query := `INSERT INTO beans(name) VALUES(?)`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Beans.Db.ExecContext(ctx, query, nw.Name)
	if err != nil {
		return err
	}

	return nil
}

func (Beans *BeansRepository) GetAll(ctx context.Context) ([]models.BeansModel, error) {

	beans := make([]models.BeansModel, 0)

	query := `SELECT id,name FROM beans WHERE is_delete=FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	rows, err := Beans.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var bean models.BeansModel

		err := rows.Scan(&bean.Id, &bean.Name)
		if err != nil {
			return nil, err
		}

		beans = append(beans, bean)
	}

	return beans, rows.Err()
}

func (Beans *BeansRepository) GetById(ctx context.Context, id int) (models.BeansModel, error) {

	var bean models.BeansModel

	query := `SELECT id,name FROM beans WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	err := Beans.Db.QueryRowContext(ctx, query, id).Scan(&bean.Id, &bean.Name)
	if err != nil {
		return models.BeansModel{}, err
	}

	return bean, nil
}

func (Beans *BeansRepository) Update(ctx context.Context, nw models.BeansModel) error {

	query := `UPDATE beans SET name = ?, is_delete = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Beans.Db.ExecContext(ctx, query, nw.Name, nw.IsDelete, nw.Id)
	if err != nil {
		return err
	}

	return nil
}

func (Beans *BeansRepository) Delete(ctx context.Context, id int) error {

	query := `UPDATE beans SET is_delete = TRUE WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	result, err := Beans.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (Beans *BeansRepository) DestroyMany(ctx context.Context) error {

	query := `DELETE FROM beans WHERE is_delete = TRUE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Beans.Db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
