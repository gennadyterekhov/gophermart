package withdrawals

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"

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

func PostHandler() http.Handler {
	return middleware.WithAuth(
		http.HandlerFunc(createWithdrawal),
		middleware.ContentTypeJson,
		middleware.Luhn,
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

func createWithdrawal(res http.ResponseWriter, req *http.Request) {
	reqDto, err := getRequestDto(req)
	if err != nil {
		logger.ZapSugarLogger.Errorln("could not getRequestDto", err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = domain.Create(req.Context(), reqDto)
	if err != nil {
		logger.ZapSugarLogger.Errorln("could not create wdr", err.Error())

		status := http.StatusInternalServerError

		if err.Error() == domain.ErrorInsufficientFunds {
			status = http.StatusPaymentRequired
		}
		http.Error(res, err.Error(), status)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func getRequestDto(req *http.Request) (*requests.Withdrawals, error) {
	requestDto := &requests.Withdrawals{
		Order: "",
		Sum:   0,
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(requestDto)
	if err != nil {
		return nil, err
	}

	if requestDto.Order == "" {
		return nil, fmt.Errorf("empty order number")
	}

	return requestDto, nil
}
