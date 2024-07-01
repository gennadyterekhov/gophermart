package login

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"

	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
)

type Controller struct {
	Service auth.Service
}

func NewController(service auth.Service) Controller {
	return Controller{
		Service: service,
	}
}

func Handler(controller *Controller) http.Handler {
	return middleware.WithoutAuth(
		http.HandlerFunc(controller.login),
		middleware.RequestContentTypeJSON,
	)
}

func (controller *Controller) login(res http.ResponseWriter, req *http.Request) {
	logger.CustomLogger.Debugln("/api/user/login handler")

	reqDto, err := getRequestDto(req)
	if err != nil {
		logger.CustomLogger.Errorln(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	resDto, err := controller.Service.Login(req.Context(), reqDto)
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == auth.ErrorWrongCredentials {
			status = http.StatusUnauthorized
		}
		logger.CustomLogger.Errorln(err.Error())

		http.Error(res, err.Error(), status)
		return
	}
	res.Header().Set("Authorization", resDto.Token)

	resBody, err := serializers.Login(resDto)
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

func getRequestDto(req *http.Request) (*requests.Login, error) {
	requestDto := &requests.Login{
		Login:    "",
		Password: "",
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(requestDto)
	if err != nil {
		return nil, err
	}

	if requestDto.Login == "" || requestDto.Password == "" {
		return nil, fmt.Errorf(auth.ErrorWrongCredentials)
	}

	return requestDto, nil
}
