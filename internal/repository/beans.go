package repository

import (
	"context"
	"database/sql"
)

type BeansRepository struct {
	Db *sql.DB
}

type BeansModel struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsDelete bool   `json:"is_delete"`
}

func (Beans *BeansRepository) Insert(ctx context.Context, nw BeansModel) error {

	query := `INSERT INTO beans(name) VALUES(?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Beans.Db.ExecContext(ctx, query, nw.Name)
	if err != nil {
		return err
	}

	return nil
}

func (Beans *BeansRepository) GetAll(ctx context.Context) ([]BeansModel, error) {

	beans := make([]BeansModel, 0)

	query := `SELECT id,name FROM beans WHERE is_delete=FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	rows, err := Beans.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var bean BeansModel

		err := rows.Scan(&bean.Id, &bean.Name)
		if err != nil {
			return nil, err
		}

		beans = append(beans, bean)
	}

	return beans, rows.Err()
}

func (Beans *BeansRepository) GetById(ctx context.Context, id int) (BeansModel, error) {

	var bean BeansModel

	query := `SELECT id,name FROM beans WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := Beans.Db.QueryRowContext(ctx, query, id).Scan(&bean.Id, &bean.Name)
	if err != nil {
		return BeansModel{}, err
	}

	return bean, nil
}

func (Beans *BeansRepository) Update(ctx context.Context, nw BeansModel) error {

	query := `UPDATE beans SET name = ?, is_delete = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Beans.Db.ExecContext(ctx, query, nw.Name, nw.IsDelete, nw.Id)
	if err != nil {
		return err
	}

	return nil
}

func (Beans *BeansRepository) Delete(ctx context.Context, id int) error {

	query := `UPDATE beans SET is_delete = TRUE WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
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

func (Beans *BeansRepository) Destroy(ctx context.Context, id int) error {

	query := `DELETE FROM beans WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Beans.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
