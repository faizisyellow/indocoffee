
package service

import (
    "context"
    "github.com/faizisyellow/indocoffee/internal/repository"
)

type RolesServices struct {
    Repository repository.Repository
}


func (Roles *RolesServices) Create(ctx context.Context) (string, error) {
 
    return "", nil
}

func (Roles *RolesServices) FindAll(ctx context.Context) ([]repository.RolesModel, error) {

    return nil, nil
}

func(Roles *RolesServices) FindById(ctx context.Context, id int) (repository.RolesModel, error) {

    return repository.RolesModel{},nil
}

func (Roles *RolesServices) Update(ctx context.Context, id int, nw repository.RolesModel) error {

    return nil
}

func (Roles *RolesServices) Delete(ctx context.Context, id int) error {

    return nil
}
