package roles

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type RolesRepository struct {
	Db *sql.DB
}

func (Roles *RolesRepository) Insert(ctx context.Context, nw models.RolesModel) error {

	qry := `INSERT INTO roles (name,level) VALUES (?,?)`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Roles.Db.ExecContext(ctx, qry, nw.Name, nw.Level)
	if err != nil {
		return err
	}

	return nil
}

func (Roles *RolesRepository) GetAll(ctx context.Context) ([]models.RolesModel, error) {

	qry := `SELECT id,name,level FROM roles WHERE is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	result, err := Roles.Db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	roles := make([]models.RolesModel, 0)

	for result.Next() {
		var role models.RolesModel
		if err := result.Scan(&role.Id, &role.Name, &role.Level); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, result.Err()
}

func (Roles *RolesRepository) GetById(ctx context.Context, id int) (models.RolesModel, error) {

	qry := `SELECT id,name,level FROM roles WHERE id  = ?  AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	result := Roles.Db.QueryRowContext(ctx, qry, id)

	var role models.RolesModel

	return role, result.Scan(&role.Id, &role.Name, &role.Level)
}

func (Roles *RolesRepository) GetByName(ctx context.Context, rolename string) (models.RolesModel, error) {

	qry := `SELECT id,name,level FROM roles WHERE name  = ?  AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	result := Roles.Db.QueryRowContext(ctx, qry, rolename)

	var role models.RolesModel

	return role, result.Scan(&role.Id, &role.Name, &role.Level)
}

func (Roles *RolesRepository) Update(ctx context.Context, nw models.RolesModel) error {

	qry := `UPDATE roles SET name = ?, level = ? WHERE id  = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Roles.Db.ExecContext(ctx, qry, nw.Name, nw.Level, nw.Id)
	if err != nil {
		return err
	}

	return nil
}

func (Roles *RolesRepository) Delete(ctx context.Context, id int) error {

	qry := `UPDATE roles SET is_delete = TRUE WHERE id  = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	result, err := Roles.Db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// if already mark as deleted
	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (Roles *RolesRepository) DestroyMany(ctx context.Context) error {

	query := `DELETE FROM roles WHERE is_delete = TRUE`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := Roles.Db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
