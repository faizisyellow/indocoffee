package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/keys"
	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/service"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware func(http.Handler) http.HandlerFunc

func NewHandlerFunc(mw ...Middleware) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		next := h
		for k := len(mw) - 1; k >= 0; k-- {
			next = mw[k](next)
			// NewHandler(m1, m2, m3)(myHandler)
			// => m1(m2(m3(myHandler)))
		}
		return next
	}
}

var (
	UsrCtx             keys.User = "user"
	ErrForbiddenAction           = errors.New("you don’t have permission for this action.")
)

func (app *Application) AuthMiddleware(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			ResponseClientError(w, r, fmt.Errorf("authorization is missing"), http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 {
			ResponseClientError(w, r, fmt.Errorf("authorization is malformed"), http.StatusUnauthorized)
			return
		}

		if parts[0] != "Bearer" {
			ResponseClientError(w, r, fmt.Errorf("authorization is malformed: authentication use Bearer"), http.StatusBadRequest)
			return
		}

		token := parts[1]

		jwtToken, err := app.Authentication.VerifyToken(token)
		if err != nil {
			ResponseClientError(w, r, err, http.StatusUnauthorized)
			return
		}

		claim, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			ResponseServerError(w, r, fmt.Errorf("error while parsing token claim type"), http.StatusInternalServerError)
			return
		}

		usrId, ok := claim["id"].(float64)
		if !ok {
			ResponseServerError(w, r, fmt.Errorf("error while parsing field claim type"), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		user, err := app.Services.UsersService.FindUserById(ctx, int(usrId))
		if err != nil {
			errService := errorService.GetError(err)
			switch errService.E {
			case service.ErrUserNotFound:
				ResponseClientError(w, r, err, http.StatusUnauthorized)
			default:
				ResponseServerError(w, r, err, http.StatusInternalServerError)
			}
			return
		}

		ctx = context.WithValue(ctx, UsrCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	}
}

func (app *Application) CheckOwnerCart(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.GetContentFromContext[*models.User](r, UsrCtx)
		if err != nil {
			ResponseServerError(w, r, err, http.StatusInternalServerError)
			return
		}

		idParam := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idParam)
		if err != nil {
			ResponseClientError(w, r, err, http.StatusBadRequest)
			return
		}

		cart, err := app.Services.CartsService.FindById(r.Context(), id)
		if err != nil {
			errService := errorService.GetError(err)
			switch errService.E {
			case service.ErrCartNotFound:
				ResponseClientError(w, r, err, http.StatusNotFound)
			default:
				ResponseServerError(w, r, err, http.StatusInternalServerError)
			}
			return
		}

		if user.Id != cart.UserId {
			ResponseClientError(w, r, ErrForbiddenAction, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
