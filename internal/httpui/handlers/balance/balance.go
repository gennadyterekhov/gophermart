package balance

import (
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/domain/balance"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"
)

type Controller struct {
	Service balance.Service
}

func NewController(service balance.Service) Controller {
	return Controller{
		Service: service,
	}
}

func Handler(controller *Controller) http.Handler {
	return middleware.WithAuth(
		http.HandlerFunc(controller.getBalance),
		middleware.ResponseContentTypeJSON,
	)
}

func (controller *Controller) getBalance(res http.ResponseWriter, req *http.Request) {
	logger.CustomLogger.Debugln(req.Method + req.RequestURI + " handler")
	resDto, err := controller.Service.GetBalanceResponse(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	resBody, err := serializers.Balance(resDto)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.CustomLogger.Debugln("returning body", string(resBody))
	_, err = res.Write(resBody)
	if err != nil {
		logger.CustomLogger.Errorln(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
