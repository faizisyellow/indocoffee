package main

import (
	"net/http"
	"strconv"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/utils"
	"github.com/go-chi/chi/v5"
)

// @Summary		Add cart
// @Description	Add new cart
// @Tags			Carts
// @Accept			json
// @Produce		json
// @Security		JWT
// @Param			payload	body		dto.CreateCartRequest	true	"Payload create new cart"
// @Success		201		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		401		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/carts [post]
func (app *Application) CreateCartsHandler(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateCartRequest
	if err := ReadHttpJson(w, r, &request); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(request); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	user, err := utils.GetContentFromContext[*models.User](r, UsrCtx)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	if err := app.Services.CartsService.Create(r.Context(), request, user.Id); err != nil {
		errService := errorService.GetError(err)
		switch errService.E {
		case service.ErrNotFoundProduct:
			ResponseClientError(w, r, err, http.StatusNotFound)
		case service.ErrConflictItemCart:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "success create cart", http.StatusCreated)
}

// @Summary		Increment cart's item
// @Description	Increment quantity's item
// @Tags			Carts
// @Accept			json
// @Produce		json
// @Security		JWT
// @Param			id	path		int	true	"cart id"
// @Success		200	{object}	main.Envelope{data=string,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		403	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/carts/{id}/increment [patch]
func (app *Application) IncrementCartsItemHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := app.Services.CartsService.IncrementItem(r.Context(), id); err != nil {
		errService := errorService.GetError(err)
		switch errService.E {
		case service.ErrCartNotFound:
			ResponseClientError(w, r, err, http.StatusNotFound)
		case service.ErrCartOverflowQuantity:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "success increment quantity item", http.StatusOK)
}

// @Summary		Decrement cart's item
// @Description	Decrement quantity's item
// @Tags			Carts
// @Accept			json
// @Produce		json
// @Security		JWT
// @Param			id	path		int	true	"cart id"
// @Success		200	{object}	main.Envelope{data=string,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		403	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/carts/{id}/decrement [patch]
func (app *Application) DecrementCartsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := app.Services.CartsService.DecrementItem(r.Context(), id); err != nil {
		errService := errorService.GetError(err)
		switch errService.E {
		case service.ErrCartNotFound:
			ResponseClientError(w, r, err, http.StatusNotFound)
		case service.ErrCartMinQuantity:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "success decrement quantity item", http.StatusOK)

}

func (app *Application) DeleteCartsHandler(w http.ResponseWriter, r *http.Request) {

}
