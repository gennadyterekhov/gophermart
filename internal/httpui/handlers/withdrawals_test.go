package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/luhn"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type withdrawalsTestSuite struct {
	suite.Suite
	tests.SuiteUsingTransactions
	tests.TestHTTPServer
	Repository repositories.Repository
}

func TestWithdrawals(t *testing.T) {
	db := storage.NewDB(tests.TestDBDSN)

	suiteInstance := &withdrawalsTestSuite{
		Repository: repositories.NewRepository(db),
	}
	suiteInstance.SetDB(db)

	suite.Run(t, suiteInstance)
}

func (suite *withdrawalsTestSuite) TestCanSendWithdrawalsRequest() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		withdrawalNewest, err := suite.Repository.AddWithdrawal(
			context.Background(),
			regDto.ID,
			"a",
			0,
			time.Time{},
		)
		assert.NoError(t, err)

		responseStatusCode, bodyAsBytes := suite.SendGet(

			"/api/user/withdrawals",
			regDto.Token,
		)

		require.Equal(t, http.StatusOK, responseStatusCode)
		responseBody := make([]responses.WithdrawalExternal, 0)
		err = json.Unmarshal(bodyAsBytes, &responseBody)
		assert.NoError(t, err)
		require.Equal(t, 1, len(responseBody))
		assert.Equal(t, withdrawalNewest.OrderNumber, (responseBody)[0].OrderNumber)
	}))
}

func (suite *withdrawalsTestSuite) Test204IfNoContent() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := suite.SendGet(

			"/api/user/withdrawals",
			regDto.Token,
		)

		assert.Equal(t, http.StatusNoContent, responseStatusCode)
	}))
}

func (suite *withdrawalsTestSuite) TestCanCreateWithdrawalsWithFloat() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		var accrual int64 = 160
		_, err := suite.Repository.AddOrder(
			context.Background(),
			"a",
			regDto.ID,
			"",
			&accrual,
			time.Time{},
		)
		require.NoError(t, err)

		orderNumber := luhn.Generate(1)
		rawJSON := fmt.Sprintf(`{"order":"%v", "sum":1.5}`, orderNumber)
		responseStatusCode := suite.SendPost(

			"/api/user/balance/withdraw",
			"application/json",
			regDto.Token,
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusOK, responseStatusCode)
	}))
}

func (suite *withdrawalsTestSuite) TestCannotCreateWithdrawalsWithIncorrectNumber() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		var accrual int64 = 10
		_, err := suite.Repository.AddOrder(
			context.Background(),
			"a",
			regDto.ID,
			"",
			&accrual,
			time.Time{},
		)
		require.NoError(t, err)

		rawJSON := `{"order":"123", "sum":1}`
		responseStatusCode := suite.SendPost(

			"/api/user/balance/withdraw",
			"application/json",
			regDto.Token,
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusUnprocessableEntity, responseStatusCode)
	}))
}

func (suite *withdrawalsTestSuite) Test402WhenNotEnoughBalance() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")

		rawJSON := `{"order":"4417123456789113", "sum":1}`
		responseStatusCode := suite.SendPost(
			"/api/user/balance/withdraw",
			"application/json",
			regDto.Token,
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusPaymentRequired, responseStatusCode)
	}))
}

func (suite *withdrawalsTestSuite) Test401IfNoToken() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := suite.SendGet(
			"/api/user/withdrawals",
			"",
		)

		assert.Equal(t, http.StatusUnauthorized, responseStatusCode)
	}))
}
