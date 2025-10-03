package main

import (
	"net/http"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

// @Summary		Create new order
// @Description	Create new order
// @Tags			Orders
// @Accept			json
// @Produce		json
// @Param			payload	body		dto.CreateOrderRequest	true	"Payload create new order"
// @Success		201			{object}	main.Envelope{data=string,error=nil}
// @Failure		400			{object}	main.Envelope{data=nil,error=string}
// @Failure		409			{object}	main.Envelope{data=nil,error=string}
// @Failure		500			{object}	main.Envelope{data=nil,error=string}
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

	if err := app.Services.OrdersService.Create(r.Context(), req, user.Id); err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrConflictBean:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "create order success", http.StatusCreated)
}
