package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	_ "github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *Application) Mux() http.Handler {

	adminRole := Admin.String()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(httprate.Limit(
		120,
		10*time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CLIENT")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Idempotency-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {

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
