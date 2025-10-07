package main

import (
	"net/http"
	"strconv"

	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/go-chi/chi/v5"
)

//	@Summary		Create coffee's form
//	@Description	Create new coffee's form
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			payload	body		dto.CreateFormRequest	true	"Form body payload"
//	@Success		201		{object}	main.Envelope{data=string,error=nil}
//	@Failure		400		{object}	main.Envelope{data=nil,error=string}
//	@Failure		401		{object}	main.Envelope{data=nil,error=string}
//	@Failure		403		{object}	main.Envelope{data=nil,error=string}
//	@Failure		409		{object}	main.Envelope{data=nil,error=string}
//	@Failure		500		{object}	main.Envelope{data=nil,error=string}
//	@Router			/forms [post]
func (app *Application) CreateFormsHandler(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateFormRequest
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	res, err := app.Services.FormsService.Create(r.Context(), req)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrConflictForm:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}

		return
	}

	ResponseSuccess(w, r, res, http.StatusCreated)

}

//	@Summary		Get coffee's form
//	@Description	Get all coffee's form
//	@Tags			Forms
//	@Produce		json
//	@Success		200	{object}	main.Envelope{data=[]dto.FormResponse,error=nil}
//	@Failure		500	{object}	main.Envelope{data=nil,error=string}
//	@Router			/forms [get]
func (app *Application) GetAllFormsHandler(w http.ResponseWriter, r *http.Request) {

	forms, err := app.Services.FormsService.FindAll(r.Context())
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	response := make([]dto.FormResponse, 0)
	for _, form := range forms {
		response = append(response, dto.FormResponse{Id: form.Id, Name: form.Name})
	}

	ResponseSuccess(w, r, response, http.StatusOK)
}

//	@Summary		Get coffee's form
//	@Description	Get coffee's form by id
//	@Tags			Forms
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		int	true	"Form id"
//	@Success		200	{object}	main.Envelope{data=dto.FormResponse,error=nil}
//	@Failure		400	{object}	main.Envelope{data=nil,error=string}
//	@Failure		401	{object}	main.Envelope{data=nil,error=string}
//	@Failure		403	{object}	main.Envelope{data=nil,error=string}
//
//	@Failure		404	{object}	main.Envelope{data=nil,error=string}
//	@Failure		500	{object}	main.Envelope{data=nil,error=string}
//	@Router			/forms/{id} [get]
func (app *Application) GetFormsHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	form, err := app.Services.FormsService.FindById(r.Context(), id)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundForm:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}

		return
	}

	response := dto.FormResponse{
		Id:   form.Id,
		Name: form.Name,
	}

	ResponseSuccess(w, r, response, http.StatusOK)
}

//	@Summary		Update coffee's form
//	@Description	Update coffee's form by id
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			id		path		int						true	"Form id"
//	@Param			payload	body		dto.UpdateFormRequest	true	"Form body payload"
//	@Success		200		{object}	main.Envelope{data=string,error=nil}
//	@Failure		400		{object}	main.Envelope{data=nil,error=string}
//	@Failure		401		{object}	main.Envelope{data=nil,error=string}
//	@Failure		403		{object}	main.Envelope{data=nil,error=string}
//	@Failure		404		{object}	main.Envelope{data=nil,error=string}
//	@Failure		409		{object}	main.Envelope{data=nil,error=string}
//	@Failure		500		{object}	main.Envelope{data=nil,error=string}
//	@Router			/forms/{id} [patch]
func (app *Application) UpdateFormsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	formId, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	var req dto.UpdateFormRequest
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	err = app.Services.FormsService.Update(r.Context(), formId, req)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundForm:
			ResponseClientError(w, r, err, http.StatusNotFound)
		case service.ErrConflictForm:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "success update form", http.StatusOK)
}

//	@Summary		Delete coffee's form
//	@Description	Delete coffee's form by id
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			id	path	int	true	"Form id"
//	@Success		204
//	@Failure		400	{object}	main.Envelope{data=nil,error=string}
//	@Failure		401	{object}	main.Envelope{data=nil,error=string}
//	@Failure		403	{object}	main.Envelope{data=nil,error=string}
//
//	@Failure		404	{object}	main.Envelope{data=nil,error=string}
//	@Failure		500	{object}	main.Envelope{data=nil,error=string}
//	@Router			/forms/{id} [delete]
func (app *Application) DeleteFormsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	formId, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	err = app.Services.FormsService.Delete(r.Context(), formId)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundForm:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)
}

//	@Summary		Delete coffee's forms
//	@Description	Delete all coffee's forms permanently
//	@Tags			Forms
//	@Security		JWT
//	@Success		204
//	@Failure		401	{object}	main.Envelope{data=nil,error=string}
//	@Failure		403	{object}	main.Envelope{data=nil,error=string}
//
//	@Failure		404	{object}	main.Envelope{data=nil,error=string}
//	@Failure		500	{object}	main.Envelope{data=nil,error=string}
//	@Router			/forms/trash [delete]
func (app *Application) TrashFormsHandler(w http.ResponseWriter, r *http.Request) {
	err := app.Services.FormsService.Remove(r.Context())
	if err != nil {
		errorValue := errorService.GetError(err)
		switch errorValue.E {
		case service.ErrNotFoundForm:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)
}
