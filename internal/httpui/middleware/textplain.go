package middleware

import (
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

func ContentTypeTextPlain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request != nil && request.Header.Get("Content-Type") != "text/plain" {
			logger.CustomLogger.Debugln("ContentTypeTextPlain middleware failed")
			response.WriteHeader(http.StatusBadRequest)
		}

		next.ServeHTTP(response, request)
	})
}
