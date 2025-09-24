package forms

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type FormsRepository struct {
	Db *sql.DB
}

func (Forms *FormsRepository) Insert(ctx context.Context, nw models.FormsModel) error {

	qry := `INSERT INTO forms (name) VALUES(?)`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Forms.Db.ExecContext(ctx, qry, nw.Name)
	if err != nil {
		return err
	}

	return nil
}

func (Forms *FormsRepository) GetAll(ctx context.Context) ([]models.FormsModel, error) {

	qry := `SELECT id,name FROM forms WHERE is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	result, err := Forms.Db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	forms := make([]models.FormsModel, 0)

	for result.Next() {
		var form models.FormsModel

		err := result.Scan(&form.Id, &form.Name)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, result.Err()
}

func (Forms *FormsRepository) GetById(ctx context.Context, id int) (models.FormsModel, error) {

	qry := `SELECT id,name FROM forms WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	result := Forms.Db.QueryRowContext(ctx, qry, id)

	var form models.FormsModel
	err := result.Scan(&form.Id, &form.Name)
	if err != nil {
		return models.FormsModel{}, err
	}

	return form, result.Err()
}

func (Forms *FormsRepository) Update(ctx context.Context, nw models.FormsModel) error {

	qry := `UPDATE forms SET name = ? WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Forms.Db.ExecContext(ctx, qry, nw.Name, nw.Id)
	if err != nil {
		return err
	}

	return nil
}

func (Forms *FormsRepository) Delete(ctx context.Context, id int) error {

	qry := `UPDATE forms SET is_delete = TRUE WHERE id = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Forms.Db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}

	return nil
}

func (Forms *FormsRepository) DestroyMany(ctx context.Context) error {

	qry := `DELETE FROM forms WHERE is_delete = TRUE `

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Forms.Db.ExecContext(ctx, qry)
	if err != nil {
		return err
	}

	return nil
}
