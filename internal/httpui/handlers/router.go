package handlers

import (
	"github.com/gennadyterekhov/gophermart/internal/client"
	"github.com/gennadyterekhov/gophermart/internal/config"
	loginDomain "github.com/gennadyterekhov/gophermart/internal/domain/auth"
	registerDomain "github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	balanceDomain "github.com/gennadyterekhov/gophermart/internal/domain/balance"
	ordersDomain "github.com/gennadyterekhov/gophermart/internal/domain/orders"
	withdrawalsDomain "github.com/gennadyterekhov/gophermart/internal/domain/withdrawals"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/balance"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/login"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/orders"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/register"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/withdrawals"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Controllers struct {
	Withdrawals withdrawals.Controller
	Orders      orders.Controller
	Balance     balance.Controller
	Register    register.Controller
	Login       login.Controller
}

type Router struct {
	Router      chi.Router
	Controllers *Controllers
}

func NewRouter(conf *config.Config, db *storage.DB) *Router {
	return &Router{
		Router:      chi.NewRouter(),
		Controllers: NewControllers(conf, db),
	}
}

func NewControllers(conf *config.Config, db *storage.DB) *Controllers {
	repo := repositories.NewRepository(db)
	balanceService := balanceDomain.NewService(repo)
	return &Controllers{
		Withdrawals: withdrawals.NewController(withdrawalsDomain.NewService(repo, balanceService)),
		Orders:      orders.NewController(ordersDomain.NewService(repo, client.NewClient(conf.AccrualURL, repo))),
		Balance:     balance.NewController(balanceService),
		Register:    register.NewController(registerDomain.NewService(repo)),
		Login:       login.NewController(loginDomain.NewService(repo)),
	}
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
