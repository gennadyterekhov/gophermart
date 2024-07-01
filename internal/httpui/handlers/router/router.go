package router

import (
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/balance"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/controllers"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/login"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/orders"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/register"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/withdrawals"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Router      chi.Router
	Controllers *controllers.Controllers
}

func NewRouter(controllers *controllers.Controllers) *Router {
	router := &Router{
		Router:      chi.NewRouter(),
		Controllers: controllers,
	}
	router.InitializeRoutes()

	return router
}

func (r Router) InitializeRoutes() {
	r.Router.Post("/api/user/login", login.Handler(&r.Controllers.Login).ServeHTTP)
	r.Router.Post("/api/user/register", register.Handler(&r.Controllers.Register).ServeHTTP)

	r.Router.Get("/api/user/balance", balance.Handler(&r.Controllers.Balance).ServeHTTP)

	r.Router.Get("/api/user/withdrawals", withdrawals.Handler(&r.Controllers.Withdrawals).ServeHTTP)
	r.Router.Post("/api/user/balance/withdraw", withdrawals.PostHandler(&r.Controllers.Withdrawals).ServeHTTP)

	r.Router.Get("/api/user/orders", orders.Handler(&r.Controllers.Orders).ServeHTTP)
	r.Router.Post("/api/user/orders", orders.PostHandler(&r.Controllers.Orders).ServeHTTP)
}
