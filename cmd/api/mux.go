package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *Application) Mux() http.Handler {

	adminRole := Admin.String()
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
			r.Get("/cart", NewHandlerFunc(app.AuthMiddleware)(app.FindUsersCartHandler))
			r.Get("/orders", NewHandlerFunc(app.AuthMiddleware)(app.FindUsersOrdersHandler))
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
			r.Post("/", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.CreateRolesHandler))
			r.Get("/", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.GetAllRolesHandler))
			r.Get("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.GetRolesHandler))
			r.Patch("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.UpdateRolesHandler))
			r.Delete("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.DeleteRolesHandler))
			r.Delete("/trash", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.TrashRolesHandler))
		})

		r.Route("/beans", func(r chi.Router) {
			r.Get("/", app.GetAllBeansHandler)
			r.Post("/", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.CreateBeansHandler))
			r.Get("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.GetBeansHandler))
			r.Patch("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.UpdateBeansHandler))
			r.Delete("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.DeleteBeansHandler))
			r.Delete("/trash", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.TrashBeansHandler))
		})

		r.Route("/forms", func(r chi.Router) {
			r.Get("/", app.GetAllFormsHandler)
			r.Post("/", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.CreateFormsHandler))
			r.Get("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.GetFormsHandler))
			r.Patch("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.UpdateFormsHandler))
			r.Delete("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.DeleteFormsHandler))
			r.Delete("/trash", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.TrashFormsHandler))
		})

		r.Route("/products", func(r chi.Router) {
			r.Get("/", app.GetProductsHandler)
			r.Get("/{id}", app.GetProductHandler)
			r.Post("/", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.CreateProductsHandler))
			r.Patch("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.UpdateProductHandler))
			r.Delete("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckAuthorization(adminRole))(app.DeleteProductHandler))
		})

		r.Route("/carts", func(r chi.Router) {
			r.Post("/", NewHandlerFunc(app.AuthMiddleware, app.OnlyActionByCustomer)(app.CreateCartsHandler))
			r.Patch("/{id}/increment", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerCart)(app.IncrementCartsItemHandler))
			r.Patch("/{id}/decrement", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerCart)(app.DecrementCartsHandler))
			r.Delete("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerCart)(app.DeleteCartsHandler))
		})

		r.Route("/orders", func(r chi.Router) {
			r.Post("/", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerCartsToOrders)(app.CreateOrdersHandler))
			r.Patch("/{id}/roast", NewHandlerFunc(app.AuthMiddleware, app.AuthorizeManageOrder)(app.ExecuteItemsHandler))
			r.Patch("/{id}/cancel", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerOrder)(app.CancelOrderHandler))
			r.Patch("/{id}/ship", NewHandlerFunc(app.AuthMiddleware, app.AuthorizeManageOrder)(app.ShipOrderHandler))
			r.Patch("/{id}/complete", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerOrder)(app.CompleteOrderHandler))
			r.Get("/{id}", NewHandlerFunc(app.AuthMiddleware, app.CheckOwnerOrder)(app.GetOrderHandler))
			r.Get("/", NewHandlerFunc(app.AuthMiddleware, app.AuthorizeManageOrder)(app.GetOrdersHandler))
		})

	})

	return r
}
