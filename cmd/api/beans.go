package main

import (
	"net/http"
	"strconv"

	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/go-chi/chi/v5"
)

// @Summary		Add Coffee bean
// @Description	Create new coffee bean
// @Tags			Beans
// @Accept			json
// @Produce		json
// @Param			payload	body		service.RequestCreateBean	true	"Payload create new bean"
// @Success		201		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/beans [post]
func (app *Application) CreateBeansHandler(w http.ResponseWriter, r *http.Request) {

	var req service.RequestCreateBean

	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	res, err := app.Services.BeansService.Create(ctx, req.Serialize())
	if err != nil {
		switch err {
		case service.ErrConflictBean:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, res, http.StatusCreated)
}

// @Summary		Get coffee beans
// @Description	Get All coffee beans
// @Tags			Beans
// @Produce		json
// @Success		200	{object}	main.Envelope{data=[]service.ResponseFindAll,error=nil}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/beans [get]
func (app *Application) GetAllBeansHandler(w http.ResponseWriter, r *http.Request) {

	res, err := app.Services.BeansService.FindAll(r.Context())
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, res, http.StatusOK)

}

// @Summary		Get coffee bean
// @Description	Get coffee bean by ID
// @Tags			Beans
// @Produce		json
// @Param			id	path		int	true	"Id coffee bean"
// @Success		200	{object}	main.Envelope{data=repository.BeansModel,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/beans/{id} [get]
func (app *Application) GetBeansHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	res, err := app.Services.BeansService.FindById(r.Context(), id)
	if err != nil {
		switch err {
		case service.ErrNotFoundBean:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, res, http.StatusOK)

}

// @Summary		Update coffee bean
// @Description	Update coffee bean by ID
// @Tags			Beans
// @Accept			json
// @Produce		json
// @Param			id		path		int							true	"Id coffee bean"
// @Param			payload	body		service.RequestUpdateBean	true	"Payload Update bean"
// @Success		200		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/beans/{id} [patch]
func (app *Application) UpdateBeansHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	var req service.RequestUpdateBean

	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	err = app.Services.BeansService.Update(r.Context(), id, req)
	if err != nil {
		switch err {
		case service.ErrNotFoundBean:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)

		}
		return
	}

	ResponseSuccess(w, r, "update bean successfully", http.StatusOK)

}

// @Summary		Delete coffee bean
// @Description	Delete coffee bean by ID
// @Tags			Beans
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Id coffee bean"
// @Success		204
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/beans/{id} [delete]
func (app *Application) DeleteBeansHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	err = app.Services.BeansService.Delete(r.Context(), id)
	if err != nil {
		switch err {
		case service.ErrNotFoundBean:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)

		}
		return
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)

}
