package main

import (
	"net/http"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

// @Summary		Get User Profile
// @Description	Get User Profile Who's log in
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		JWT
// @Success		200	{object}	main.Envelope{data=dto.GetUsersProfileResponse,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/users/profile [get]
func (app *Application) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {

	user, err := utils.GetContentFromContext[*models.User](r, UsrCtx)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	response := dto.GetUsersProfileResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  *user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	ResponseSuccess(w, r, response, http.StatusOK)
}

// @Summary		Delete User Account
// @Description	Delete User Account
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		JWT
// @Success		204
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/users/delete [delete]
func (app *Application) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {

	user, err := utils.GetContentFromContext[*models.User](r, UsrCtx)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = app.Services.UsersService.DeleteAccount(r.Context(), user.Id)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)
}

// @Summary		Get User's cart
// @Description	Get User's cart
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		JWT
// @Success		200	{object}	main.Envelope{data=dto.GetUsersCartResponse}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/users/cart [get]
func (app *Application) FindUsersCartHandler(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetContentFromContext[*models.User](r, UsrCtx)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	userWithCarts, err := app.Services.UsersService.FindUsersCart(r.Context(), user.Id)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, userWithCarts, http.StatusOK)
}
