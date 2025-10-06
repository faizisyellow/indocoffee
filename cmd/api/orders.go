package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/utils"
	"github.com/go-chi/chi/v5"
)

// @Summary		Create new order
// @Description	Create new order
// @Tags			Orders
// @Accept			json
// @Produce		json
// @Param			X-Idempotency-Key	header	string					true	"unique identifier each order"
// @Param			payload				body	dto.CreateOrderRequest	true	"Payload create new order"
// @Security		JWT
// @Success		201	{object}	main.Envelope{data=string,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		403	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		409	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/orders [post]
func (app *Application) CreateOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	user, err := utils.GetContentFromContext[*models.User](r, UsrCtx)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	idempotencyHeader := r.Header.Get("X-Idempotency-Key")
	if idempotencyHeader == "" {
		ResponseClientError(w, r, errors.New("require an idempotency key"), http.StatusBadRequest)
	}

	if len(idempotencyHeader) != 64 {
		ResponseClientError(w, r, errors.New("idempotency key is malformed: require 32 bytes string"), http.StatusBadRequest)
	}

	if err := app.Services.OrdersService.Create(r.Context(), idempotencyHeader, req, user.Id); err != nil {
		if strings.Contains(err.Error(), "phone validation") {
			ResponseClientError(w, r, err, http.StatusBadRequest)
			return
		}

		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrOrdersItemEmpty:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		case service.ErrOrdersQuantityIssue:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		case service.ErrCartMinQuantity:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		case service.ErrOrdersConflict:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "create order success", http.StatusCreated)
}

// @Summary		Roast Order
// @Description	Perform the roasting process for an order
// @Tags			Orders
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"Order id"
// @Security		JWT
// @Success		200	{object}	main.Envelope{data=string,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		403	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/orders/{id}/roast [patch]
func (app *Application) ExecuteItemsHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.Services.OrdersService.ExecuteItems(r.Context(), chi.URLParam(r, "id")); err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrOrdersNotFound:
			ResponseClientError(w, r, err, http.StatusNotFound)
		case service.ErrOrdersInvalidStatus:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "success update order to be roasting", http.StatusOK)
}
