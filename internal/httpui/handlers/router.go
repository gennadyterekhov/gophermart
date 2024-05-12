package handlers

import (
	"github.com/go-chi/chi/v5"
)

func GetRouter() chi.Router {
	router := chi.NewRouter()
	registerRoutes(router)

	return router
}

func registerRoutes(router chi.Router) {
	router.Head("/", HeadHandler)

	router.Get("/api/user/withdrawals", TempHandler().ServeHTTP)
	router.Get("/api/user/balance", TempHandler().ServeHTTP)
	router.Get("/api/user/orders", TempHandler().ServeHTTP)
	router.Post("/api/user/balance/withdraw", TempHandler().ServeHTTP)
	router.Post("/api/user/orders", TempHandler().ServeHTTP)
	router.Post("/api/user/login", TempHandler().ServeHTTP)
	router.Post("/api/user/register", TempHandler().ServeHTTP)
}
