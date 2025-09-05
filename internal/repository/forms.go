package repository

import (
	"context"
	"database/sql"
)

type FormsRepository struct {
	Db *sql.DB
}

type FormsModel struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsDelete string `json:"is_delete"`
}

func (Forms *FormsRepository) Insert(ctx context.Context, nw FormsModel) error {

	qry := `INSERT INTO forms (name) VALUES(?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Forms.Db.ExecContext(ctx, qry, nw.Name)
	if err != nil {
		return err
	}

	return nil
}

func (Forms *FormsRepository) GetAll(ctx context.Context) ([]FormsModel, error) {

	qry := `SELECT id,name FROM forms WHERE is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	result, err := Forms.Db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	forms := make([]FormsModel, 0)

	for result.Next() {
		var form FormsModel

		err := result.Scan(&form.Id, &form.Name)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, result.Err()
}

func (Forms *FormsRepository) GetById(ctx context.Context, id int) (FormsModel, error) {

	qry := `SELECT id,name FROM forms WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	result := Forms.Db.QueryRowContext(ctx, qry, id)

	var form FormsModel
	err := result.Scan(&form.Id, &form.Name)
	if err != nil {
		return FormsModel{}, err
	}

	return form, result.Err()
}

func (Forms *FormsRepository) Update(ctx context.Context, nw FormsModel) error {

	qry := `UPDATE forms SET name = ?, is_delete = ? WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Forms.Db.ExecContext(ctx, qry, nw.Name, nw.IsDelete, nw.Id)
	if err != nil {
		return err
	}

	return nil
}

func (Forms *FormsRepository) Delete(ctx context.Context, id int) error {

	qry := `UPDATE forms SET is_delete = TRUE WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Forms.Db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}

	return nil
}

func (Forms *FormsRepository) Destroy(ctx context.Context, id int) error {

	qry := `DELETE FROM forms WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Forms.Db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}

	return nil
}
