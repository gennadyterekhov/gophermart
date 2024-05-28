package middleware

import (
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request != nil && request.Header.Get("Content-Type") != "application/json" {
			logger.ZapSugarLogger.Debugln("ContentTypeJSON middleware failed")
			response.WriteHeader(http.StatusBadRequest)
		}

		next.ServeHTTP(response, request)
	})
}
