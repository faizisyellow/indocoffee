package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *Application) Mux() http.Handler {

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Route("/health", func(r chi.Router) {
			r.Get("/", app.HealthHandler)
		})

		r.Route("/uploadthing", func(r chi.Router) {
			// for callback upload before upload file
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("done"))
			})
		})

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
			r.Delete("/trash", app.TrashRolesHandler)
		})

		r.Route("/beans", func(r chi.Router) {
			r.Post("/", app.CreateBeansHandler)
			r.Get("/", app.GetAllBeansHandler)
			r.Get("/{id}", app.GetBeansHandler)
			r.Patch("/{id}", app.UpdateBeansHandler)
			r.Delete("/{id}", app.DeleteBeansHandler)
			r.Delete("/trash", app.TrashBeansHandler)
		})

		r.Route("/forms", func(r chi.Router) {
			r.Post("/", app.CreateFormsHandler)
			r.Get("/", app.GetAllFormsHandler)
			r.Get("/{id}", app.GetFormsHandler)
			r.Patch("/{id}", app.UpdateFormsHandler)
			r.Delete("/{id}", app.DeleteFormsHandler)
			r.Delete("/trash", app.TrashFormsHandler)
		})

		r.Route("/products", func(r chi.Router) {
			r.Post("/", app.CreateProductsHandler)
			r.Get("/{id}", app.GetProductHandler)
			r.Get("/", app.GetProductsHandler)
			r.Patch("/{id}", app.UpdateProductHandler)
			r.Delete("/{id}", app.DeleteProductHandler)
		})

		r.Route("/carts", func(r chi.Router) {
			r.Post("/", NewHandlerFunc(app.AuthMiddleware)(app.CreateCartsHandler))
			r.Patch("/{id}/increment", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerCart)(app.IncrementCartsItemHandler))
			r.Patch("/{id}/decrement", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerCart)(app.DecrementCartsHandler))
			r.Delete("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerCart)(app.DeleteCartsHandler))
		})
	})

	return r
}
