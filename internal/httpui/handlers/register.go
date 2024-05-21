package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"

	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
)

func RegisterHandler() http.Handler {
	return middleware.WithoutAuth(
		http.HandlerFunc(register),
		middleware.ContentTypeJson,
	)
}

func register(res http.ResponseWriter, req *http.Request) {
	var err error
	reqDto, err := getRequestDto(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	resDto, err := auth.Register(req.Context(), reqDto)
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == auth.ErrorNotUniqueLogin {
			status = http.StatusConflict
		}
		http.Error(res, err.Error(), status)
		return
	}
	resBody := serializers.Register(resDto)

	_, err = res.Write([]byte(resBody))
	if err != nil {
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
