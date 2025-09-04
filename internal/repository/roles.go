package repository

import (
	"context"
	"database/sql"
)

type RolesRepository struct {
	Db *sql.DB
}

type RolesModel struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func (Roles *RolesRepository) Insert(ctx context.Context, nw RolesModel) error {

	qry := `INSERT INTO roles (name,level) VALUES (?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Roles.Db.ExecContext(ctx, qry, nw.Name, nw.Level)
	if err != nil {
		return err
	}

	return nil
}

func (Roles *RolesRepository) GetAll(ctx context.Context) ([]RolesModel, error) {

	qry := `SELECT id,name,level FROM roles WHERE is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	result, err := Roles.Db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	roles := make([]RolesModel, 0)

	for result.Next() {
		var role RolesModel
		if err := result.Scan(&role.Id, &role.Name, &role.Level); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, result.Err()
}

func (Roles *RolesRepository) GetById(ctx context.Context, id int) (RolesModel, error) {

	qry := `SELECT id,name,level FROM roles WHERE id  = ?  AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	result := Roles.Db.QueryRowContext(ctx, qry, id)

	var role RolesModel

	return role, result.Scan(&role.Id, &role.Name, &role.Level)
}

func (Roles *RolesRepository) Update(ctx context.Context, nw RolesModel) error {

	qry := `UPDATE roles SET name = ?, level = ? WHERE id  = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Roles.Db.ExecContext(ctx, qry, nw.Name, nw.Level, nw.Id)
	if err != nil {
		return err
	}

	return nil
}

func (Roles *RolesRepository) Delete(ctx context.Context, id int) error {

	qry := `UPDATE roles SET is_delete = TRUE WHERE id  = ? AND is_delete = FALSE`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
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

func (Roles *RolesRepository) Destroy(ctx context.Context, id int) error {

	query := `DELETE FROM roles WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := Roles.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
