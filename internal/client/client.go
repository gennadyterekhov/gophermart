package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/gennadyterekhov/gophermart/internal/config"
)

type AccrualClientResponse struct {
	CorrectResponse         *CorrectResponse
	TooManyRequestsResponse *TooManyRequestsResponse
}

type CorrectResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type TooManyRequestsResponse struct {
	RetryAfter        int64
	RequestsPerMinute int64
}

const (
	ErrorNoContent       = "order is not registered"
	ErrorInternal        = "internal server error"
	ErrorUnknownResponse = "unknown response"
)

// SendOrderToAccrual sends orders and returns status. if already exists, does not create new order
func SendOrderToAccrual(number string) (*AccrualClientResponse, error) {
	var err error
	path := fmt.Sprintf("/api/orders/%v", number)
	url := config.ServerConfig.AccrualURL + path
	var client *resty.Client = resty.New()
	response, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	statusCode := response.StatusCode()

	responseDto := &AccrualClientResponse{}
	if statusCode == http.StatusOK {
		correctResponse, err := processSuccessfulResponse(response)
		if err != nil {
			return nil, err
		}
		responseDto.CorrectResponse = correctResponse
		return responseDto, nil
	}
	if statusCode == http.StatusTooManyRequests {
		tooManyRequestsResponse, err := process409Response(response)
		if err != nil {
			return nil, err
		}
		responseDto.TooManyRequestsResponse = tooManyRequestsResponse
		return responseDto, nil
	}
	if statusCode == http.StatusNoContent {
		return responseDto, fmt.Errorf(ErrorNoContent)
	}
	if statusCode == http.StatusInternalServerError {
		return responseDto, fmt.Errorf(ErrorInternal)
	}

	return responseDto, fmt.Errorf(ErrorUnknownResponse)
}

func processSuccessfulResponse(response *resty.Response) (*CorrectResponse, error) {
	responseDto := &CorrectResponse{}

	err := json.Unmarshal(response.Body(), responseDto)
	if err != nil {
		return nil, err
	}

	// TODO update statuses https://github.com/gennadyterekhov/gophermart/issues/14
	return responseDto, nil
}

func process409Response(response *resty.Response) (*TooManyRequestsResponse, error) {
	responseDto := &TooManyRequestsResponse{}

	retryAfterRaw := response.Header().Get("Retry-After")
	retryAfter, err := strconv.ParseInt(retryAfterRaw, 10, 64)
	if err != nil {
		return nil, err
	}
	responseDto.RetryAfter = retryAfter

	bodyAsString := string(response.Body()) //  No more than N requests per minute allowed
	bodyAsString = strings.Replace(bodyAsString, "No more than ", "", 1)
	bodyAsString = strings.Replace(bodyAsString, " requests per minute allowed", "", 1)

	requestsPerMinute, err := strconv.ParseInt(strings.TrimSpace(bodyAsString), 10, 64)
	if err != nil {
		return nil, err
	}
	responseDto.RequestsPerMinute = requestsPerMinute

	return responseDto, nil
}