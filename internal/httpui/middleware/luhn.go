package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/luhn"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

func Luhn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var reqBody []byte
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			logger.CustomLogger.Errorln("could not read body", err.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		var number string
		if req.Header.Get("Content-Type") == "text/plain" {
			number = getNumberFromTextBody(reqBody)
		} else {
			number, err = getNumberFromJSONBody(reqBody)
		}

		if err != nil {
			logger.CustomLogger.Errorln(err.Error())
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		err = luhn.Validate(number)
		if err != nil {
			logger.CustomLogger.Errorln(err.Error())
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func getNumberFromTextBody(reqBody []byte) string {
	return string(reqBody)
}

func getNumberFromJSONBody(reqBody []byte) (string, error) {
	var err error
	type jsonWithOrder struct {
		Order string `json:"order,omitempty"`
	}
	jsonWithOrderInstance := jsonWithOrder{}

	err = json.Unmarshal(reqBody, &jsonWithOrderInstance)
	if err != nil {
		return "", err
	}

	if jsonWithOrderInstance.Order == "" {
		return "", fmt.Errorf("order field is empty")
	}

	return jsonWithOrderInstance.Order, nil
}
