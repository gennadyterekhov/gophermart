package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/with_server"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/luhn"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type withdrawalsTestSuite struct {
	with_server.BaseSuiteWithServer
}

func TestWithdrawals(t *testing.T) {
	suiteInstance := &withdrawalsTestSuite{}
	with_server.InitBaseSuiteWithServer(suiteInstance)

	suite.Run(t, suiteInstance)
}

func (suite *withdrawalsTestSuite) TestCanSendWithdrawalsRequest() {
	regDto := suite.RegisterForTest("a", "a")
	withdrawalNewest, err := suite.Repository.AddWithdrawal(
		context.Background(),
		regDto.ID,
		"a",
		0,
		time.Time{},
	)
	assert.NoError(suite.T(), err)

	responseStatusCode, bodyAsBytes := suite.SendGet(

		"/api/user/withdrawals",
		regDto.Token,
	)

	require.Equal(suite.T(), http.StatusOK, responseStatusCode)
	responseBody := make([]responses.WithdrawalExternal, 0)
	err = json.Unmarshal(bodyAsBytes, &responseBody)
	assert.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, len(responseBody))
	assert.Equal(suite.T(), withdrawalNewest.OrderNumber, (responseBody)[0].OrderNumber)
}

func (suite *withdrawalsTestSuite) Test204IfNoContent() {
	regDto := suite.RegisterForTest("a", "a")

	responseStatusCode, _ := suite.SendGet(

		"/api/user/withdrawals",
		regDto.Token,
	)

	assert.Equal(suite.T(), http.StatusNoContent, responseStatusCode)
}

func (suite *withdrawalsTestSuite) TestCanCreateWithdrawalsWithFloat() {
	regDto := suite.RegisterForTest("a", "a")
	var accrual int64 = 160
	_, err := suite.Repository.AddOrder(
		context.Background(),
		"a",
		regDto.ID,
		"",
		&accrual,
		time.Time{},
	)
	require.NoError(suite.T(), err)

	orderNumber := luhn.Generate(1)
	rawJSON := fmt.Sprintf(`{"order":"%v", "sum":1.5}`, orderNumber)
	responseStatusCode := suite.SendPost(

		"/api/user/balance/withdraw",
		"application/json",
		regDto.Token,
		bytes.NewBuffer([]byte(rawJSON)),
	)

	assert.Equal(suite.T(), http.StatusOK, responseStatusCode)
}

func (suite *withdrawalsTestSuite) TestCannotCreateWithdrawalsWithIncorrectNumber() {
	regDto := suite.RegisterForTest("a", "a")
	var accrual int64 = 10
	_, err := suite.Repository.AddOrder(
		context.Background(),
		"a",
		regDto.ID,
		"",
		&accrual,
		time.Time{},
	)
	require.NoError(suite.T(), err)

	rawJSON := `{"order":"123", "sum":1}`
	responseStatusCode := suite.SendPost(

		"/api/user/balance/withdraw",
		"application/json",
		regDto.Token,
		bytes.NewBuffer([]byte(rawJSON)),
	)

	assert.Equal(suite.T(), http.StatusUnprocessableEntity, responseStatusCode)
}

func (suite *withdrawalsTestSuite) Test402WhenNotEnoughBalance() {
	regDto := suite.RegisterForTest("a", "a")

	rawJSON := `{"order":"4417123456789113", "sum":1}`
	responseStatusCode := suite.SendPost(
		"/api/user/balance/withdraw",
		"application/json",
		regDto.Token,
		bytes.NewBuffer([]byte(rawJSON)),
	)

	assert.Equal(suite.T(), http.StatusPaymentRequired, responseStatusCode)
}

func (suite *withdrawalsTestSuite) Test401IfNoToken() {
	suite.RegisterForTest("a", "a")

	responseStatusCode, _ := suite.SendGet(
		"/api/user/withdrawals",
		"",
	)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseStatusCode)
}
