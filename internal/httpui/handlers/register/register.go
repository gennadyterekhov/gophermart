package register

import (
	"encoding/json"
	"net/http"

	domain "github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/httpui/serializers"
)

func Handler() http.Handler {
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
		return
	}

	resDto, err := domain.Register(req.Context(), reqDto)
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == domain.ErrorNotUniqueLogin {
			status = http.StatusConflict
		}
		http.Error(res, err.Error(), status)
		return
	}
	resBody, err := serializers.Register(resDto)
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
