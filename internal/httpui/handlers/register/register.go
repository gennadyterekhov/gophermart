package register

import (
	"encoding/json"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	domain "github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"
)

type Controller struct {
	Service domain.Service
}

func NewController(service domain.Service) Controller {
	return Controller{
		Service: service,
	}
}

func Handler(controller *Controller) http.Handler {
	return middleware.WithoutAuth(
		http.HandlerFunc(controller.register),
		middleware.RequestContentTypeJSON,
	)
}

func (controller *Controller) register(res http.ResponseWriter, req *http.Request) {
	logger.CustomLogger.Debugln("/api/user/register handler")

	var err error
	reqDto, err := getRequestDto(req)
	if err != nil {
		logger.CustomLogger.Errorln(err.Error())

		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	resDto, err := controller.Service.Register(req.Context(), reqDto)
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == domain.ErrorNotUniqueLogin {
			status = http.StatusConflict
		}
		logger.CustomLogger.Errorln(err.Error())

		http.Error(res, err.Error(), status)
		return
	}
	res.Header().Set("Authorization", resDto.Token)

	resBody, err := serializers.Register(resDto)
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

func getRequestDto(req *http.Request) (*requests.Register, error) {
	requestDto := &requests.Register{
		Login:    "",
		Password: "",
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(requestDto)
	if err != nil {
		return nil, err
	}

	return requestDto, nil
}
