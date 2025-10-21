package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
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
// @Success		201	{object}	main.Envelope{data=dto.CreateOrderResponse,error=nil}
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
		return
	}

	if valid := utils.ValidateIdempotencyKey(idempotencyHeader); !valid {
		ResponseClientError(w, r, errors.New("idempotency key malformed: must be UUIDv4"), http.StatusBadRequest)
		return
	}

	newOrderId, err := app.Services.OrdersService.Create(r.Context(), idempotencyHeader, req, user.Id)
	if err != nil {
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

	ResponseSuccess(w, r, dto.CreateOrderResponse{Id: newOrderId}, http.StatusCreated)
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

// @Summary		Cancel Order
// @Description	Perform cancel for an order
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
// @Router			/orders/{id}/cancel [patch]
func (app *Application) CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.Services.OrdersService.CancelOrder(r.Context(), chi.URLParam(r, "id")); err != nil {
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

	ResponseSuccess(w, r, "success update order to be cancelled", http.StatusOK)
}

// @Summary		Ship Order
// @Description	Perform ship for an order
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
// @Router			/orders/{id}/ship [patch]
func (app *Application) ShipOrderHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.Services.OrdersService.ShipOrder(r.Context(), chi.URLParam(r, "id")); err != nil {
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

	ResponseSuccess(w, r, "success update order to be shipped", http.StatusOK)
}

// @Summary		Complete Order
// @Description	Perform complete for an order
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
// @Router			/orders/{id}/complete [patch]
func (app *Application) CompleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.Services.OrdersService.CompleteOrder(r.Context(), chi.URLParam(r, "id")); err != nil {
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

	ResponseSuccess(w, r, "success order complete", http.StatusOK)
}

// @Summary		Get Order
// @Description	Get Order by id
// @Tags			Orders
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"Order id"
// @Security		JWT
// @Success		200	{object}	main.Envelope{data=dto.GetOrderResponse,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		403	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/orders/{id} [get]
func (app *Application) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		ResponseClientError(w, r, errors.New("path param is empty"), http.StatusBadRequest)
	}

	order, err := app.Services.OrdersService.FindById(r.Context(), id)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrOrdersNotFound:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, dto.GetOrderResponse{
		Id:                     order.Id,
		CustomerName:           order.CustomerName,
		CustomerEmail:          order.CustomerEmail,
		Status:                 order.Status,
		Items:                  order.Items,
		TotalPrice:             order.TotalPrice,
		PhoneNumber:            order.PhoneNumber,
		AlternativePhoneNumber: order.AlternativePhoneNumber,
		Street:                 order.Street,
		City:                   order.City,
		CreatedAt:              order.CreatedAt,
	}, http.StatusOK)
}

// @Summary		Get Orders
// @Description	Get All Orders
// @Tags			Orders
// @Accept			json
// @Produce		json
// @Security		JWT
// @Param			status	query		string	false	"status order"
// @Param			sort	query		string	false	"sort order by created asc(oldest) | desc(latest)"
// @Param			limit	query		string	false	"limit each page"
// @Param			offset	query		string	false	"skip rows"
// @Success		200		{object}	main.Envelope{data=[]dto.GetOrderResponse,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		401		{object}	main.Envelope{data=nil,error=string}
// @Failure		403		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/orders [get]
func (app *Application) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	requestQuery := repository.QueryOrders{
		Limit:  r.URL.Query().Get("limit"),
		Offset: r.URL.Query().Get("offset"),
		Sort:   r.URL.Query().Get("sort"),
		Status: r.URL.Query().Get("status"),
	}

	paginateOrders, err := repository.PaginatedOrdersQuery{Limit: 10, Sort: "asc"}.Parse(requestQuery)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(paginateOrders); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	orders, err := app.Services.OrdersService.FindOrders(r.Context(), paginateOrders)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	var response []dto.GetOrderResponse
	for _, order := range orders {
		response = append(response, dto.GetOrderResponse{
			Id:                     order.Id,
			CustomerName:           order.CustomerName,
			CustomerEmail:          order.CustomerEmail,
			Status:                 order.Status,
			Items:                  order.Items,
			TotalPrice:             order.TotalPrice,
			PhoneNumber:            order.PhoneNumber,
			AlternativePhoneNumber: order.AlternativePhoneNumber,
			Street:                 order.Street,
			City:                   order.City,
			CreatedAt:              order.CreatedAt,
		})
	}

	ResponseSuccess(w, r, response, http.StatusOK)

}
