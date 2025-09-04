package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/repository"
	serviceParser "github.com/faizisyellow/indocoffee/internal/service/parser"
)

type RolesServices struct {
	Repository repository.Repository
}

var (
	ErrConflictRole = errors.New("roles: role already exist")
	ErrInternalRole = errors.New("roles: encountered an internal error")
	ErrNotFoundRole = errors.New("roles: no such as role")
)

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required,min=4"`
}

func (Cr CreateRoleRequest) Serialize() CreateRoleRequest {

	Cr.Name = strings.ToLower(Cr.Name)
	return Cr
}

func (Roles *RolesServices) Create(ctx context.Context, req CreateRoleRequest) (string, error) {

	var newRole repository.RolesModel
	newRole.Name = req.Name

	err := Roles.Repository.Roles.Insert(ctx, newRole)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", ErrConflictRole
		}

		return "", ErrInternalRole
	}

	return "success create new role", nil
}

type ResponseRolesFindAll struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (rr *ResponseRolesFindAll) ParseDTO(data any) error {
	switch v := data.(type) {
	case repository.RolesModel:
		rr.Id = v.Id
		rr.Name = v.Name
	default:
		return ErrInternalRole
	}

	return nil
}

func (Roles *RolesServices) FindAll(ctx context.Context) ([]ResponseRolesFindAll, error) {

	roles, err := Roles.Repository.Roles.GetAll(ctx)
	if err != nil {
		return nil, ErrInternalRole
	}

	response := make([]ResponseRolesFindAll, 0)

	for _, role := range roles {
		res := new(ResponseRolesFindAll)
		err := serviceParser.Parse(res, role)
		if err != nil {
			return nil, err
		}
		response = append(response, *res)
	}

	return response, nil
}

type ResponseRolesById struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (rr *ResponseRolesById) ParseDTO(data any) error {
	switch v := data.(type) {
	case repository.RolesModel:
		rr.Id = v.Id
		rr.Name = v.Name
	default:
		return ErrInternalRole
	}

	return nil
}

func (Roles *RolesServices) FindById(ctx context.Context, id int) (ResponseRolesById, error) {

	role, err := Roles.Repository.Roles.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ResponseRolesById{}, ErrNotFoundRole
		default:
			return ResponseRolesById{}, ErrInternalRole
		}
	}

	response := new(ResponseRolesById)
	err = serviceParser.Parse(response, role)
	if err != nil {
		return ResponseRolesById{}, ErrInternalRole
	}

	return *response, nil
}

type RequestUpdateRole struct {
	Name string `json:"name" validate:"min=4"`
}

func (ru RequestUpdateRole) Serialize() RequestUpdateRole {

	ru.Name = strings.ToLower(ru.Name)
	return ru
}

func UpdateRolePayload(req RequestUpdateRole, ltsRole repository.RolesModel) repository.RolesModel {

	ltsRole.Name = req.Name

	return ltsRole
}

func (Roles *RolesServices) Update(ctx context.Context, id int, nw repository.RolesModel) error {

	err := Roles.Repository.Roles.Update(ctx, nw)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return ErrConflictRole
		}

		return ErrInternalRole
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
		return ErrInternalRole
	}

	return nil
}

func (Roles *RolesServices) Remove(ctx context.Context, id int) error {

	role, err := Roles.FindById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFoundRole
		default:
			return ErrInternalRole
		}
	}

	err = Roles.Repository.Roles.Destroy(ctx, role.Id)
	if err != nil {
		return ErrInternalRole
	}

	return nil
}
