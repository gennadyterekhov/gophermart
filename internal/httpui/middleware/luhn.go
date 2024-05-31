package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

const (
	asciiZero = 48
	asciiTen  = 57
)

func Luhn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var reqBody []byte
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			logger.ZapSugarLogger.Errorln("could not read body", err.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		var number string
		if req.Header.Get("Content-Type") == "text/plain" {
			number, err = getNumberFromTextBody(reqBody)
		} else {
			number, err = getNumberFromJSONBody(reqBody)
		}

		if err != nil {
			logger.ZapSugarLogger.Error(err.Error())
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		err = isOk(number)
		if err != nil {
			logger.ZapSugarLogger.Error(err.Error())
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func getNumberFromTextBody(reqBody []byte) (string, error) {
	return string(reqBody), nil
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

func isOk(number string) error {
	parity := len(number) % 2
	sum, err := calculateLuhnSum(number, parity)
	if err != nil {
		return err
	}

	// If the total modulo 10 is not equal to 0, then the number is invalid.
	if sum%10 != 0 {
		return fmt.Errorf("invalid number")
	}

	return nil
}

func calculateLuhnSum(number string, parity int) (int64, error) {
	var sum int64
	for i, d := range number {
		if d < asciiZero || d > asciiTen {
			return 0, fmt.Errorf("invalid digit")
		}

		d = d - asciiZero
		// Double the value of every second digit.
		if i%2 == parity {
			d *= 2
			// If the result of this doubling operation is greater than 9.
			if d > 9 {
				// The same final result can be found by subtracting 9 from that result.
				d -= 9
			}
		}

		// Take the sum of all the digits.
		sum += int64(d)
	}

	return sum, nil
}
