package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/with_server"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type balanceTestSuite struct {
	with_server.BaseSuiteWithServer
}

func (suite *balanceTestSuite) SetupSuite() {
	with_server.InitBaseSuiteWithServer(suite)
}

func TestBalance(t *testing.T) {
	suite.Run(t, new(balanceTestSuite))
}

func (suite *balanceTestSuite) TestCanSendBalanceRequest() {
	regDto := suite.RegisterForTest("a", "a")
	_, err := suite.Repository.AddWithdrawal(context.Background(), regDto.ID, "", 100, time.Time{})
	assert.NoError(suite.T(), err)

	var startBalance int64 = 1000 // cents => $10
	_, err = suite.Repository.AddOrder(context.Background(), "", regDto.ID, "", &startBalance, time.Time{})
	assert.NoError(suite.T(), err)

	responseStatusCode, bodyAsBytes := suite.SendGet(
		"/api/user/balance",
		regDto.Token,
	)

	require.Equal(suite.T(), http.StatusOK, responseStatusCode)
	responseBody := &responses.BalanceExternal{}
	err = json.Unmarshal(bodyAsBytes, responseBody)
	assert.NoError(suite.T(), err)
	require.Equal(suite.T(), float64(1), responseBody.Withdrawn)
	assert.Equal(suite.T(), float64(9), responseBody.Current)
}

func (suite *balanceTestSuite) TestBalance401IfNoToken() {
	suite.RegisterForTest("a", "a")

	responseStatusCode, _ := suite.SendGet(
		"/api/user/balance",
		"",
	)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseStatusCode)
}
