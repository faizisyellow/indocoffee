package main

import (
	"net/http"
	"strconv"

	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/go-chi/chi/v5"
)

// @Summary		Add Coffee bean
// @Description	Create new coffee bean
// @Tags			Beans
// @Accept			json
// @Produce		json
// @Param			payload	body		dto.CreateBeanRequest	true	"Payload create new bean"
// @Success		201		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/beans [post]
func (app *Application) CreateBeansHandler(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateBeanRequest
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	res, err := app.Services.BeansService.Create(r.Context(), req)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
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
// @Success		200	{object}	main.Envelope{data=[]dto.BeanResponse,error=nil}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/beans [get]
func (app *Application) GetAllBeansHandler(w http.ResponseWriter, r *http.Request) {

	beans, err := app.Services.BeansService.FindAll(r.Context())
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	response := make([]dto.BeanResponse, 0)
	for _, bean := range beans {
		response = append(response, dto.BeanResponse{Id: bean.Id, Name: bean.Name})
	}

	ResponseSuccess(w, r, response, http.StatusOK)

}

// @Summary		Get coffee bean
// @Description	Get coffee bean by ID
// @Tags			Beans
// @Produce		json
// @Param			id	path		int	true	"Id coffee bean"
// @Success		200	{object}	main.Envelope{data=dto.BeanResponse,error=nil}
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

	bean, err := app.Services.BeansService.FindById(r.Context(), id)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundBean:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	response := dto.BeanResponse{
		Id:   bean.Id,
		Name: bean.Name,
	}

	ResponseSuccess(w, r, response, http.StatusOK)

}

// @Summary		Update coffee bean
// @Description	Update coffee bean by ID
// @Tags			Beans
// @Accept			json
// @Produce		json
// @Param			id		path		int						true	"Id coffee bean"
// @Param			payload	body		dto.UpdateBeanRequest	true	"Payload Update bean"
// @Success		200		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/beans/{id} [patch]
func (app *Application) UpdateBeansHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	beanId, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	var req dto.UpdateBeanRequest
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	err = app.Services.BeansService.Update(r.Context(), beanId, req)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundBean:
			ResponseClientError(w, r, err, http.StatusNotFound)
		case service.ErrConflictBean:
			ResponseClientError(w, r, err, http.StatusConflict)
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
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundBean:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)

}

// @Summary		Delete coffee's beans
// @Description	Delete all coffee's beans permanently
// @Tags			Beans
// @Success		204
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/beans/trash [delete]
func (app *Application) TrashBeansHandler(w http.ResponseWriter, r *http.Request) {

	err := app.Services.BeansService.Remove(r.Context())
	if err != nil {
		errorValue := errorService.GetError(err)
		switch errorValue.E {
		case service.ErrNotFoundBean:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)
}
