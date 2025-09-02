package main

import (
	"fmt"
	"net/http"

	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// @Summary		Sign Up Account
// @Description	Sign Up New Account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			payload	body		service.RegisterRequest	true	"Payload to Sign Up"
// @Success		201		{object}	main.Envelope{data=service.RegisterResponse,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/authentication/sign-up [post]
func (app *Application) SignUpHandler(w http.ResponseWriter, r *http.Request) {

	var req service.RegisterRequest

	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	res, err := app.Services.UsersService.RegisterAccount(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrUserAlreadyExist:
			ResponseClientError(w, r, fmt.Errorf("email not available"), http.StatusConflict)
		case utils.ErrInvalidPasswordSignature:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, res, http.StatusCreated)

}

// @Summary		Activate Account
// @Description	Activate  New Account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			token	path		string	true	"Token Inivitation to Activate Account"
// @Success		200		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/authentication/activation/{token} [post]
func (app *Application) ActivateAccountHandler(w http.ResponseWriter, r *http.Request) {

	req := chi.URLParam(r, "token")

	if req == "" {
		ResponseClientError(w, r, fmt.Errorf("token invition is required"), http.StatusBadRequest)
		return
	}

	err := app.Services.UsersService.ActivateAccount(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrTokenInvitationNotFound, service.ErrUserRegisteredNotFound:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "user activated successfully", http.StatusOK)

}

type LoginResponse struct {
	Token string `json:"token"`
}

// @Summary		Sign in Account
// @Description	Sign in  New Account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			payload	body		service.LoginRequest	true	"Token to Sign in Account"
// @Success		200		{object}	main.Envelope{data=main.LoginResponse,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/authentication/sign-in [post]
func (app *Application) SignInHandler(w http.ResponseWriter, r *http.Request) {

	var req service.LoginRequest

	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	user, err := app.Services.UsersService.Login(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			ResponseClientError(w, r, err, http.StatusNotFound)
		case service.ErrUserNotActivated:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		case bcrypt.ErrMismatchedHashAndPassword:
			ResponseClientError(w, r, fmt.Errorf("email or password incorrect"), http.StatusBadRequest)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}

		return
	}

	claims := jwt.MapClaims{
		"iss": app.JwtAuth.Iss,
		"sub": app.JwtAuth.Sub,
		"exp": app.JwtAuth.Exp,
		"id":  user.Id,
	}

	token, err := app.Authentication.GenerateToken(claims)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(w, r, LoginResponse{Token: token}, http.StatusOK)
}
