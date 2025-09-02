package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *Application) Mux() http.Handler {

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.Get("/profile", NewHandlerFunc(app.AuthMiddleware)(app.GetUserProfileHandler))
			r.Delete("/delete", NewHandlerFunc(app.AuthMiddleware)(app.DeleteAccountHandler))
		})

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(app.SwaggerUrl)))

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/sign-up", app.SignUpHandler)
			r.Post("/activation/{token}", app.ActivateAccountHandler)
			r.Post("/sign-in", app.SignInHandler)
		})

		r.Route("/roles", func(r chi.Router) {
			r.Post("/", app.CreateRolesHandler)
			r.Get("/", app.GetAllRolesHandler)
			r.Get("/{id}", app.GetRolesHandler)
			r.Patch("/{id}", app.UpdateRolesHandler)
			r.Delete("/{id}", app.DeleteRolesHandler)
		})
	})

	return r
}
