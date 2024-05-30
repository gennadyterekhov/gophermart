package balance

import (
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/domain/balance"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"
)

func Handler() http.Handler {
	return middleware.WithAuth(
		http.HandlerFunc(getBalance),
		middleware.ContentTypeJSON,
	)
}

func getBalance(res http.ResponseWriter, req *http.Request) {
	resDto, err := balance.GetBalanceResponse(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	resBody, err := serializers.Balance(resDto)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = res.Write(resBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
