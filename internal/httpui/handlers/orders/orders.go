package orders

import (
	"io"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"

	domain "github.com/gennadyterekhov/gophermart/internal/domain/orders"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"
	"github.com/gennadyterekhov/gophermart/internal/logger"
)

func Handler() http.Handler {
	return middleware.WithAuth(
		http.HandlerFunc(orders),
		middleware.ResponseContentTypeJSON,
	)
}

func PostHandler() http.Handler {
	return middleware.WithAuth(
		http.HandlerFunc(sendOrderToProcessing),
		middleware.ContentTypeTextPlain,
		middleware.Luhn,
	)
}

func orders(res http.ResponseWriter, req *http.Request) {
	logger.CustomLogger.Debugln(req.Method + req.RequestURI + " handler")

	resDto, err := domain.GetAll(req.Context())
	if err != nil {
		logger.CustomLogger.Errorln(err.Error())

		status := http.StatusInternalServerError

		if err.Error() == domain.ErrorNoContent {
			status = http.StatusNoContent
		}
		http.Error(res, err.Error(), status)
		return
	}

	resBody, err := serializers.Orders(resDto)
	if err != nil {
		logger.CustomLogger.Errorln(err.Error())

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

func sendOrderToProcessing(res http.ResponseWriter, req *http.Request) {
	logger.CustomLogger.Debugln(req.Method + req.RequestURI + " handler")

	reqDto, err := getRequestDto(req)
	if err != nil {
		logger.CustomLogger.Errorln("could not getRequestDto", err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = domain.Create(req.Context(), reqDto)
	if err != nil {

		if err.Error() == domain.ErrorNumberAlreadyUploaded {
			res.WriteHeader(http.StatusOK)
			return
		}
		logger.CustomLogger.Errorln(err.Error())

		if err.Error() == domain.ErrorNumberAlreadyUploadedByAnotherUser {
			http.Error(res, err.Error(), http.StatusConflict)
			return
		}

		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func getRequestDto(req *http.Request) (*requests.Orders, error) {
	requestDto := &requests.Orders{
		Number: "",
	}

	readBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()

	requestDto.Number = string(readBytes)

	return requestDto, nil
}
