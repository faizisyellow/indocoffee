package main

import (
	"net/http"
	"strconv"

	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/go-chi/chi/v5"
)

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

	res, err := app.Services.BeansService.Create(ctx, req)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, res, http.StatusOK)
}

func (app *Application) GetAllBeansHandler(w http.ResponseWriter, r *http.Request) {

	res, err := app.Services.BeansService.FindAll(r.Context())
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, res, http.StatusOK)

}

func (app *Application) GetBeansHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	res, err := app.Services.BeansService.FindById(r.Context(), id)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, res, http.StatusOK)

}

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

	var updateBean repository.BeansModel
	updateBean.Name = req.Name

	err = app.Services.BeansService.Update(r.Context(), id, updateBean)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, "update bean successfully", http.StatusOK)

}

func (app *Application) DeleteBeansHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	err = app.Services.BeansService.Delete(r.Context(), id)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)

}
