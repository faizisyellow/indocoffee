package main

import (
	"net/http"

	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

// @Summary		Get User Profile
// @Description	Get User Profile Who's log in
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		JWT
// @Success		200	{object}	main.Envelope{data=repository.UserModel,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/users/profile [get]
func (app *Application) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {

	user, err := utils.GetContentFromContext[*repository.UserModel](r, UsrCtx)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, user, http.StatusOK)
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

	user, err := utils.GetContentFromContext[*repository.UserModel](r, UsrCtx)
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
