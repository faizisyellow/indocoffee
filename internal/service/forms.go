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

type FormsServices struct {
	Repository repository.Repository
}

var (
	ErrConflictForm = errors.New("forms: form already exist")
	ErrInternalForm = errors.New("forms: encountered an internal error")
	ErrNotFoundForm = errors.New("forms: no such as form")
)

func (Forms *FormsServices) Create(ctx context.Context, req dto.CreateFormRequest) (string, error) {

	newForm := repository.FormsModel{
		Name: req.Name,
	}

	err := Forms.Repository.Forms.Insert(ctx, newForm)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", errorService.New(ErrConflictForm, err)
		}

		return "", errorService.New(ErrInternalForm, err)
	}

	return "success create new form", nil
}

func (Forms *FormsServices) FindAll(ctx context.Context) ([]repository.FormsModel, error) {

	forms, err := Forms.Repository.Forms.GetAll(ctx)
	if err != nil {
		return nil, errorService.New(ErrInternalForm, err)
	}

	return forms, nil
}

func (Forms *FormsServices) FindById(ctx context.Context, id int) (repository.FormsModel, error) {

	form, err := Forms.Repository.Forms.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return repository.FormsModel{}, errorService.New(ErrNotFoundForm, err)
		default:
			return repository.FormsModel{}, errorService.New(ErrInternalForm, err)
		}
	}

	return form, nil
}

func (Forms *FormsServices) Update(ctx context.Context, id int, req dto.UpdateFormRequest) error {

	form, err := Forms.FindById(ctx, id)
	if err != nil {
		return err
	}
	form.Name = req.Name

	err = Forms.Repository.Forms.Update(ctx, form)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return errorService.New(ErrConflictForm, err)
		}

		return errorService.New(ErrInternalForm, err)
	}

	return nil
}

func (Forms *FormsServices) Delete(ctx context.Context, id int) error {

	form, err := Forms.FindById(ctx, id)
	if err != nil {
		return err
	}

	err = Forms.Repository.Forms.Delete(ctx, form.Id)
	if err != nil {
		return errorService.New(ErrInternalForm, err)
	}

	return nil
}

func removeFormsWithConcurrent(repo repository.Repository, ctx context.Context) error {

	return nil
}

func (Forms *FormsServices) Remove(ctx context.Context) error {

	return removeRolesWithConcurrent(Forms.Repository, ctx)
}
