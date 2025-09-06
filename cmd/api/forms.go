package main

import (
	"net/http"
	"strconv"

	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/go-chi/chi/v5"
)

// @Summary		Create coffee's form
// @Description	Create new coffee's form
// @Tags			Forms
// @Accept			json
// @Produce		json
// @Param			payload	body		service.CreateFormRequest	true	"Form body payload"
// @Success		201		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/forms [post]
func (app *Application) CreateFormsHandler(w http.ResponseWriter, r *http.Request) {

	var req service.CreateFormRequest
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	res, err := app.Services.FormsService.Create(r.Context(), req.Serialize())
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

// @Summary		Get coffee's form
// @Description	Get all coffee's form
// @Tags			Forms
// @Produce		json
// @Success		200	{object}	main.Envelope{data=[]service.ResponseFormsFindAll,error=nil}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/forms [get]
func (app *Application) GetAllFormsHandler(w http.ResponseWriter, r *http.Request) {

	res, err := app.Services.FormsService.FindAll(r.Context())
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, res, http.StatusOK)
}

// @Summary		Get coffee's form
// @Description	Get coffee's form by id
// @Tags			Forms
// @Produce		json
// @Param			id	path		int	true	"Form id"
// @Success		200	{object}	main.Envelope{data=[]service.ResponseFormsFindAll,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/forms/{id} [get]
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

	ResponseSuccess(w, r, form, http.StatusOK)
}

// @Summary		Update coffee's form
// @Description	Update coffee's form by id
// @Tags			Forms
// @Accept			json
// @Produce		json
// @Param			id		path		int							true	"Form id"
// @Param			payload	body		service.UpdateFormRequest	true	"Form body payload"
// @Success		200		{object}	main.Envelope{data=[]service.ResponseFormsFindAll,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/forms/{id} [patch]
func (app *Application) UpdateFormsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	formReq, err := app.Services.FormsService.FindById(ctx, id)
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

	var req service.UpdateFormRequest
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	var form repository.FormsModel
	form.Id = formReq.Id
	form.Name = formReq.Name

	form = service.UpdateFormPayload(req, form)
	err = app.Services.FormsService.Update(ctx, form.Id, form)
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

	ResponseSuccess(w, r, "success update form", http.StatusOK)
}

// @Summary		Delete coffee's form
// @Description	Delete coffee's form by id
// @Tags			Forms
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Form id"
// @Success		204
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/forms/{id} [delete]
func (app *Application) DeleteFormsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	err = app.Services.FormsService.Delete(r.Context(), id)
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

// @Summary		Delete coffee's forms
// @Description	Delete all coffee's forms permanently
// @Tags			Forms
// @Success		204
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/forms [delete]
func (app *Application) TrashFormsHandler(w http.ResponseWriter, r *http.Request) {

}
