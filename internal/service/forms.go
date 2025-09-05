package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/repository"
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

type CreateFormRequest struct {
	Name string `json:"name" validate:"required,min=4"`
}

func (cr CreateFormRequest) Serialize() CreateFormRequest {

	cr.Name = strings.ToLower(cr.Name)
	return cr
}

func (Forms *FormsServices) Create(ctx context.Context, req CreateFormRequest) (string, error) {

	var newForm repository.FormsModel
	newForm.Name = req.Name

	err := Forms.Repository.Forms.Insert(ctx, newForm)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", errorService.New(ErrConflictForm.Error(), err.Error())
		}

		return "", errorService.New(ErrInternalForm.Error(), err.Error())
	}

	return "success create new form", nil
}

type ResponseFormsFindAll struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (rf *ResponseFormsFindAll) ParseDTO(data any) error {

	switch v := data.(type) {
	case repository.FormsModel:
		rf.Id = v.Id
		rf.Name = v.Name
	default:
		return errors.New("parse responseFormsFiendAll: unknown type")
	}

	return nil
}

func (Forms *FormsServices) FindAll(ctx context.Context) ([]ResponseFormsFindAll, error) {

	forms, err := Forms.Repository.Forms.GetAll(ctx)
	if err != nil {
		return nil, errorService.New(ErrInternalForm.Error(), err.Error())
	}

	response := make([]ResponseFormsFindAll, 0)

	for _, form := range forms {
		res := new(ResponseFormsFindAll)
		err := res.ParseDTO(form)
		if err != nil {
			return nil, errorService.New(ErrInternalForm.Error(), err.Error())
		}
		response = append(response, *res)
	}

	return response, nil
}

type ResponseFormsById struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (rfb *ResponseFormsById) ParseDTO(data any) error {

	switch v := data.(type) {
	case repository.FormsModel:
		rfb.Id = v.Id
		rfb.Name = v.Name
	default:
		return errors.New("parse responseFormsById: unknown type")
	}

	return nil
}

func (Forms *FormsServices) FindById(ctx context.Context, id int) (ResponseFormsById, error) {

	form, err := Forms.Repository.Forms.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ResponseFormsById{}, errorService.New(ErrNotFoundForm.Error(), err.Error())
		default:
			return ResponseFormsById{}, errorService.New(ErrInternalForm.Error(), err.Error())
		}
	}

	var response ResponseFormsById
	err = response.ParseDTO(form)
	if err != nil {
		return ResponseFormsById{}, errorService.New(ErrInternalForm.Error(), err.Error())
	}

	return response, nil
}

type UpdateFormRequest struct {
	Name string `json:"name" validate:"required,min=4"`
}

func UpdateFormPayload(req UpdateFormRequest, lts repository.FormsModel) repository.FormsModel {

	lts.Name = req.Name
	return lts
}

func (Forms *FormsServices) Update(ctx context.Context, id int, nw repository.FormsModel) error {

	err := Forms.Repository.Forms.Update(ctx, nw)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return errorService.New(ErrConflictForm.Error(), err.Error())
		}

		return errorService.New(ErrInternalForm.Error(), err.Error())
	}

	return nil
}

func (Forms *FormsServices) Delete(ctx context.Context, id int) error {

	form, err := Forms.FindById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errorService.New(ErrNotFoundForm.Error(), err.Error())
		default:
			return errorService.New(ErrInternalForm.Error(), err.Error())
		}
	}

	err = Forms.Repository.Forms.Delete(ctx, form.Id)
	if err != nil {
		return errorService.New(ErrInternalForm.Error(), err.Error())
	}

	return nil
}

func (Forms *FormsServices) Remove(ctx context.Context, id int) error {

	form, err := Forms.FindById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errorService.New(ErrNotFoundForm.Error(), err.Error())
		default:
			return errorService.New(ErrInternalForm.Error(), err.Error())
		}
	}

	err = Forms.Repository.Forms.Destroy(ctx, form.Id)
	if err != nil {
		return errorService.New(ErrInternalForm.Error(), err.Error())
	}

	return nil
}
