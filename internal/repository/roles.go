
package repository

import (
	"database/sql"
    "context"
)

type RolesRepository struct {
	Db *sql.DB
}

type RolesModel struct {

}

func (Roles *RolesRepository) Insert(ctx context.Context) error {

    return nil
}

func (Roles *RolesRepository) GetAll(ctx context.Context) ([]RolesModel, error) {

    return nil, nil
}

func (Roles *RolesRepository) GetById(ctx context.Context, id int) (RolesModel, error) {

    return RolesModel{}, nil
}

func (Roles *RolesRepository) Update(ctx context.Context, nw RolesModel) error {

    return  nil
}

func (Roles *RolesRepository) Delete(ctx context.Context, id int) error {

    return  nil
}
