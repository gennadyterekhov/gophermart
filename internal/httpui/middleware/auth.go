package middleware

import (
	"context"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
)

const (
	ContextUserIDKey = "user_id"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		id, _, err := auth.GetIDAndLoginFromToken(authHeader)
		if err != nil {
			logger.ZapSugarLogger.Error(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), ContextUserIDKey, id)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
