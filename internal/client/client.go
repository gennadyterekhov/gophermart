package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"

	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/go-resty/resty/v2"
)

type AccrualClientResponse struct {
	CorrectResponse         *CorrectResponse
	NoContentResponse       *NoContentResponse
	TooManyRequestsResponse *TooManyRequestsResponse
}

type CorrectResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type NoContentResponse struct {
	Status string `json:"status"`
}

type TooManyRequestsResponse struct {
	RetryAfter        int64
	RequestsPerMinute int64
}

type AccrualClient struct {
	AccrualURL string
	Repository repositories.RepositoryInterface
}

func NewClient(url string, repo repositories.RepositoryInterface) AccrualClient {
	return AccrualClient{
		AccrualURL: url,
		Repository: repo,
	}
}

const (
	ErrorNoContent       = "order is not registered"
	ErrorInternal        = "internal server error"
	ErrorUnknownResponse = "unknown response"
)

func (ac *AccrualClient) GetStatus(number string) (*AccrualClientResponse, error) {
	var err error
	path := fmt.Sprintf("/api/orders/%v", number)
	url := ac.AccrualURL + path
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
		correctResponse, err := processNoContentResponse(response)
		if err != nil {
			return nil, err
		}
		responseDto.NoContentResponse = correctResponse
		return responseDto, nil // fmt.Errorf(ErrorNoContent)
	}
	if statusCode == http.StatusInternalServerError {
		return responseDto, fmt.Errorf(ErrorInternal)
	}

	return responseDto, fmt.Errorf(ErrorUnknownResponse)
}

func (ac *AccrualClient) RegisterOrderInAccrual(number string) (int, error) {
	bodyBytes := []byte(`
			{
				"order": "` + number + `",
				"goods": [
					{
						"description": "Стиральная машинка LG",
						"price": 47399.99
					}
				]
			}
		`)
	var client *resty.Client = resty.New()
	url := ac.AccrualURL + "/api/orders"
	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyBytes)

	resp, err := req.Post(url)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode(), nil
}

func processSuccessfulResponse(response *resty.Response) (*CorrectResponse, error) {
	responseDto := &CorrectResponse{}

	err := json.Unmarshal(response.Body(), responseDto)
	if err != nil {
		return nil, err
	}

	return responseDto, nil
}

func processNoContentResponse(response *resty.Response) (*NoContentResponse, error) {
	responseDto := &NoContentResponse{}
	responseDto.Status = order.New

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
