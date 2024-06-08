package middleware

import (
	"context"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/token"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

const (
	ContextUserIDKey ContextStorageKey = "user_id"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		id, _, err := token.GetIDAndLoginFromToken(authHeader)
		if err != nil {
			logger.CustomLogger.Errorln(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), ContextUserIDKey, id)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
