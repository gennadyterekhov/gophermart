package withdrawals

import (
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	domain "github.com/gennadyterekhov/gophermart/internal/domain/withdrawals"
	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"

	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
)

func Handler() http.Handler {
	return middleware.WithAuth(
		http.HandlerFunc(withdrawals),
		middleware.ContentTypeJson,
	)
}

func withdrawals(res http.ResponseWriter, req *http.Request) {
	resDto, err := domain.GetAll(req.Context())
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())

		status := http.StatusInternalServerError

		if err.Error() == domain.ErrorNoContent {
			status = http.StatusNoContent
		}
		http.Error(res, err.Error(), status)
		return
	}

	resBody, err := serializers.Withdrawals(resDto)
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())

		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = res.Write(resBody)
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())

		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
