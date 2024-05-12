package middleware

import (
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request != nil && request.Header.Get("Content-Type") == "application/json" {
			logger.ZapSugarLogger.Debugln("ContentType set to json ")
			response.Header().Set("Content-Type", "application/json")
		} else {
			logger.ZapSugarLogger.Debugln("ContentType set to text/html")
			response.Header().Set("Content-Type", "text/html")
		}

		next.ServeHTTP(response, request)
	})
}
