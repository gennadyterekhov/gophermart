package login

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"

	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
)

func Handler() http.Handler {
	return middleware.WithoutAuth(
		http.HandlerFunc(login),
		middleware.ContentTypeJSON,
	)
}

func login(res http.ResponseWriter, req *http.Request) {
	reqDto, err := getRequestDto(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	resDto, err := auth.Login(req.Context(), reqDto)
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == auth.ErrorWrongCredentials {
			status = http.StatusUnauthorized
		}
		http.Error(res, err.Error(), status)
		return
	}
	resBody, err := serializers.Login(resDto)
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
