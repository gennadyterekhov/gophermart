package handlers

import (
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
)

func temp(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "not implemented", http.StatusNotFound)
}

func TempHandler() http.Handler {
	return middleware.WithAuth(
		http.HandlerFunc(temp),
	)
}
