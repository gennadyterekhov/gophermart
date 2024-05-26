package handlers

import (
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/login"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/register"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/withdrawals"
	"github.com/go-chi/chi/v5"
)

func GetRouter() chi.Router {
	router := chi.NewRouter()
	registerRoutes(router)

	return router
}

func registerRoutes(router chi.Router) {
	router.Head("/", HeadHandler)

	router.Post("/api/user/login", login.Handler().ServeHTTP)
	router.Post("/api/user/register", register.Handler().ServeHTTP)

	router.Get("/api/user/balance", TempHandler().ServeHTTP)

	router.Get("/api/user/withdrawals", withdrawals.Handler().ServeHTTP)
	router.Post("/api/user/balance/withdraw", TempHandler().ServeHTTP)

	router.Get("/api/user/orders", TempHandler().ServeHTTP)
	router.Post("/api/user/orders", TempHandler().ServeHTTP)
}
