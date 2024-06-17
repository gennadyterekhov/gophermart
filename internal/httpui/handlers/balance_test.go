package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type balanceTestSuite struct {
	suite.Suite
	tests.SuiteUsingTransactions
	tests.TestHTTPServer
	Repository repositories.Repository
}

func TestBalance(t *testing.T) {
	db := storage.NewDB(tests.TestDBDSN)

	suiteInstance := &balanceTestSuite{
		Repository: repositories.NewRepository(db),
	}
	suiteInstance.SetDB(db)

	suite.Run(t, suiteInstance)
}

func (suite *balanceTestSuite) TestCanSendBalanceRequest() {
	run := suite.UsingTransactions()
	suite.T().Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		_, err := suite.Repository.AddWithdrawal(context.Background(), regDto.ID, "", 100, time.Time{})
		assert.NoError(t, err)

		var startBalance int64 = 1000 // cents => $10
		_, err = suite.Repository.AddOrder(context.Background(), "", regDto.ID, "", &startBalance, time.Time{})
		assert.NoError(t, err)

		responseStatusCode, bodyAsBytes := suite.SendGet(
			"/api/user/balance",
			regDto.Token,
		)

		require.Equal(t, http.StatusOK, responseStatusCode)
		responseBody := &responses.BalanceExternal{}
		err = json.Unmarshal(bodyAsBytes, responseBody)
		assert.NoError(t, err)
		require.Equal(t, float64(1), responseBody.Withdrawn)
		assert.Equal(t, float64(9), responseBody.Current)
	}))
}

func (suite *balanceTestSuite) TestBalance401IfNoToken() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := suite.SendGet(
			"/api/user/balance",
			"",
		)

		assert.Equal(t, http.StatusUnauthorized, responseStatusCode)
	}))
}
