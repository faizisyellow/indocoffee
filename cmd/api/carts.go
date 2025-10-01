package main

import (
	"net/http"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/utils"
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

func (app *Application) DeleteCartsHandler(w http.ResponseWriter, r *http.Request) {

}
