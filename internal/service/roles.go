package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/repository"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
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
	Name  string `json:"name" validate:"required,min=4"`
	Level int    `json:"level" validate:"required,min=1"`
}

func (Cr CreateRoleRequest) Serialize() CreateRoleRequest {

	Cr.Name = strings.ToLower(Cr.Name)
	return Cr
}

func (Roles *RolesServices) Create(ctx context.Context, req CreateRoleRequest) (string, error) {

	var newRole repository.RolesModel
	newRole.Name = req.Name
	newRole.Level = req.Level

	err := Roles.Repository.Roles.Insert(ctx, newRole)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", errorService.New(ErrConflictRole, err)
		}

		return "", errorService.New(ErrInternalRole, err)
	}

	return "success create new role", nil
}

type ResponseRolesFindAll struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func (rr *ResponseRolesFindAll) ParseDTO(data any) error {
	switch v := data.(type) {
	case repository.RolesModel:
		rr.Id = v.Id
		rr.Name = v.Name
		rr.Level = v.Level
	default:
		return errors.New("parse responseRolesFindAll: unknown type")
	}

	return nil
}

func (Roles *RolesServices) FindAll(ctx context.Context) ([]ResponseRolesFindAll, error) {

	roles, err := Roles.Repository.Roles.GetAll(ctx)
	if err != nil {
		return nil, errorService.New(ErrInternalRole, err)
	}

	response := make([]ResponseRolesFindAll, 0)

	for _, role := range roles {
		res := new(ResponseRolesFindAll)
		err := serviceParser.Parse(res, role)
		if err != nil {
			return nil, errorService.New(ErrInternalRole, err)
		}
		response = append(response, *res)
	}

	return response, nil
}

type ResponseRolesById struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func (rr *ResponseRolesById) ParseDTO(data any) error {
	switch v := data.(type) {
	case repository.RolesModel:
		rr.Id = v.Id
		rr.Name = v.Name
		rr.Level = v.Level
	default:
		return errors.New("parse responseRolesById: unknown type")
	}

	return nil
}

func (Roles *RolesServices) FindById(ctx context.Context, id int) (ResponseRolesById, error) {

	role, err := Roles.Repository.Roles.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ResponseRolesById{}, errorService.New(ErrNotFoundRole, err)
		default:
			return ResponseRolesById{}, errorService.New(ErrInternalRole, err)
		}
	}

	response := new(ResponseRolesById)
	err = serviceParser.Parse(response, role)
	if err != nil {
		return ResponseRolesById{}, errorService.New(ErrInternalRole, err)
	}

	return *response, nil
}

type RequestUpdateRole struct {
	Name  string `json:"name" validate:"min=4"`
	Level *int   `json:"level" validate:"omitempty,min=1"`
}

func (ru RequestUpdateRole) Serialize() RequestUpdateRole {

	ru.Name = strings.ToLower(ru.Name)
	return ru
}

func UpdateRolePayload(req RequestUpdateRole, ltsRole repository.RolesModel) repository.RolesModel {

	ltsRole.Name = req.Name
	if req.Level != nil {
		ltsRole.Level = *req.Level
	}

	return ltsRole
}

func (Roles *RolesServices) Update(ctx context.Context, id int, nw repository.RolesModel) error {

	err := Roles.Repository.Roles.Update(ctx, nw)
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

func removeRolesWithConcurrent(repo repository.Repository, ctx context.Context) error {

	return nil
}

func (Roles *RolesServices) Remove(ctx context.Context) error {

	// can switch services
	return removeRolesWithConcurrent(Roles.Repository, ctx)
}
