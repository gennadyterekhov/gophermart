package handlers

import (
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
	reqDto, errCode, err := getRequestDto(req)
	if err != nil {
		http.Error(res, err.Error(), errCode)
	}

	resDto, err := auth.Register(req.Context(), reqDto)

	resBody := serializers.Register(resDto)
	if err != nil {
		http.Error(res, err.Error(), 400)
		return
	}

	_, err = res.Write([]byte(resBody))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func getRequestDto(req *http.Request) (*requests.Register, int, error) {
	validate(req)
	return nil, 0, nil
}

func validate(req *http.Request) bool {
	return true
}
