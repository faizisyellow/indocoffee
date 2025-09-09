package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
)

type RolesServices struct {
	Repository repository.Repository
}

var (
	ErrConflictRole = errors.New("roles: role already exist")
	ErrInternalRole = errors.New("roles: encountered an internal error")
	ErrNotFoundRole = errors.New("roles: no such as role")
	ErrUpdateRole   = errors.New("roles: fields not specify")
)

func (Roles *RolesServices) Create(ctx context.Context, req dto.CreateRoleRequest) (string, error) {

	newRole := repository.RolesModel{
		Name:  req.Name,
		Level: req.Level,
	}

	err := Roles.Repository.Roles.Insert(ctx, newRole)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", errorService.New(ErrConflictRole, err)
		}

		return "", errorService.New(ErrInternalRole, err)
	}

	return "success create new role", nil
}

func (Roles *RolesServices) FindAll(ctx context.Context) ([]repository.RolesModel, error) {

	roles, err := Roles.Repository.Roles.GetAll(ctx)
	if err != nil {
		return nil, errorService.New(ErrInternalRole, err)
	}

	return roles, nil
}

func (Roles *RolesServices) FindById(ctx context.Context, id int) (repository.RolesModel, error) {

	role, err := Roles.Repository.Roles.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return repository.RolesModel{}, errorService.New(ErrNotFoundRole, err)
		default:
			return repository.RolesModel{}, errorService.New(ErrInternalRole, err)
		}
	}

	return role, nil
}

func (Roles *RolesServices) Update(ctx context.Context, id int, req dto.UpdateRoleRequest) error {

	existingRole, err := Roles.FindById(ctx, id)
	if err != nil {
		return err
	}

	if req.Name == nil && req.Level == nil {
		return errorService.New(ErrUpdateRole, errors.New("roles: update fields request not specify"))
	}

	if req.Name != nil {
		existingRole.Name = *req.Name
	}

	if req.Level != nil {
		existingRole.Level = *req.Level
	}

	err = Roles.Repository.Roles.Update(ctx, existingRole)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return errorService.New(ErrConflictRole, err)
		}

		return errorService.New(ErrInternalRole, err)
	}

	return nil
}

func (Roles *RolesServices) Delete(ctx context.Context, id int) error {

	role, err := Roles.FindById(ctx, id)
	if err != nil {
		return err
	}

	err = Roles.Repository.Roles.Delete(ctx, role.Id)
	if err != nil {
		return errorService.New(ErrInternalRole, err)
	}

	return nil
}

func (Roles *RolesServices) Remove(ctx context.Context) error {

	err := Roles.Repository.Roles.DestroyMany(ctx)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errorService.New(ErrNotFoundRole, err)
		default:
			return errorService.New(ErrConflictRole, err)
		}
	}

	return nil
}
