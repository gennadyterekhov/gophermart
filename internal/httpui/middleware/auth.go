package middleware

import (
	"context"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
)

const ContextUserKey = "login"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		login, err := auth.GetLoginFromToken(authHeader)
		if err != nil {
			logger.ZapSugarLogger.Error(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), ContextUserKey, login)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
