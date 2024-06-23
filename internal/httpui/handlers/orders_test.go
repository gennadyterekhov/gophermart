package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/with_server"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ordersTestSuite struct {
	with_server.BaseSuiteWithServer
}

func newSuite() *ordersTestSuite {
	suiteInstance := &ordersTestSuite{}
	with_server.InitBaseSuiteWithServer(suiteInstance)

	return suiteInstance
}

func TestOrders(t *testing.T) {
	suite.Run(t, newSuite())
}

func (suite *ordersTestSuite) TestCanSendOrdersRequest() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		var err error
		regDto := suite.RegisterForTest("a", "a")
		_, err = suite.Repository.AddOrder(
			context.Background(),
			"1",
			regDto.ID,
			"st1", nil,
			time.Time{},
		)
		assert.NoError(t, err)

		var tenDollarsOrThousandCents int64 = 1000
		_, err = suite.Repository.AddOrder(
			context.Background(),
			"2",
			regDto.ID,
			"st2", &tenDollarsOrThousandCents,
			time.Time{}.AddDate(1, 0, 0),
		)
		assert.NoError(t, err)

		responseStatusCode, bodyAsBytes := suite.SendGet(
			"/api/user/orders",
			regDto.Token,
		)

		require.Equal(t, http.StatusOK, responseStatusCode)
		responseBody := make([]order.OrderFloats, 0)

		err = json.Unmarshal(bodyAsBytes, &responseBody)
		assert.NoError(t, err)
		require.Equal(t, 2, len(responseBody))
		assert.Nil(t, (responseBody)[0].Accrual)
		assert.Equal(t, float64(10), *(responseBody)[1].Accrual)
	}))
}

func (suite *ordersTestSuite) TestOrders204IfNoContent() {
	run := suite.UsingTransactions()
	suite.T().Run("", run(func(t *testing.T) {
		regDto := suite.RegisterForTest("a", "a")

		responseStatusCode, _ := suite.SendGet(
			"/api/user/orders",
			regDto.Token,
		)

		assert.Equal(t, http.StatusNoContent, responseStatusCode)
	}))
}

func (suite *ordersTestSuite) TestOrders401IfNoToken() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		suite.RegisterForTest("a", "a")

		responseStatusCode, _ := suite.SendGet(
			"/api/user/orders",
			"",
		)

		assert.Equal(t, http.StatusUnauthorized, responseStatusCode)
	}))
}

func (suite *ordersTestSuite) Test200IfAlreadyUploaded() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		var err error
		regDto := suite.RegisterForTest("a", "a")
		_, err = suite.Repository.AddOrder(
			context.Background(),
			"12345678903",
			regDto.ID,
			"st1", nil,
			time.Time{},
		)
		assert.NoError(t, err)

		responseStatusCode := suite.SendPost(
			"/api/user/orders",
			"text/plain",
			regDto.Token,
			bytes.NewBuffer([]byte("12345678903")),
		)

		require.Equal(t, http.StatusOK, responseStatusCode)
	}))
}

func (suite *ordersTestSuite) Test409IfAlreadyUploadedByAnotherUser() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		var err error
		anotherUser := suite.RegisterForTest("another", "a")

		_, err = suite.Repository.AddOrder(
			context.Background(),
			"12345678903",
			anotherUser.ID,
			"st1", nil,
			time.Time{},
		)
		assert.NoError(t, err)

		regDto := suite.RegisterForTest("a", "a")
		responseStatusCode := suite.SendPost(
			"/api/user/orders",
			"text/plain",
			regDto.Token,
			bytes.NewBuffer([]byte("12345678903")),
		)

		require.Equal(t, http.StatusConflict, responseStatusCode)
	}))
}

func (suite *ordersTestSuite) Test422IfInvalidNumber() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		var _ error
		regDto := suite.RegisterForTest("a", "a")

		responseStatusCode := suite.SendPost(

			"/api/user/orders",
			"text/plain",
			regDto.Token,
			bytes.NewBuffer([]byte("1234567890")),
		)

		require.Equal(t, http.StatusUnprocessableEntity, responseStatusCode)
	}))
}
