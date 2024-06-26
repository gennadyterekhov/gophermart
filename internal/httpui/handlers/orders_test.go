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
	var err error
	regDto := suite.RegisterForTest("a", "a")
	_, err = suite.Repository.AddOrder(
		context.Background(),
		"1",
		regDto.ID,
		"st1", nil,
		time.Time{},
	)
	assert.NoError(suite.T(), err)

	var tenDollarsOrThousandCents int64 = 1000
	_, err = suite.Repository.AddOrder(
		context.Background(),
		"2",
		regDto.ID,
		"st2", &tenDollarsOrThousandCents,
		time.Time{}.AddDate(1, 0, 0),
	)
	assert.NoError(suite.T(), err)

	responseStatusCode, bodyAsBytes := suite.SendGet(
		"/api/user/orders",
		regDto.Token,
	)

	require.Equal(suite.T(), http.StatusOK, responseStatusCode)
	responseBody := make([]order.OrderFloats, 0)

	err = json.Unmarshal(bodyAsBytes, &responseBody)
	assert.NoError(suite.T(), err)
	require.Equal(suite.T(), 2, len(responseBody))
	assert.Nil(suite.T(), (responseBody)[0].Accrual)
	assert.Equal(suite.T(), float64(10), *(responseBody)[1].Accrual)
}

func (suite *ordersTestSuite) TestOrders204IfNoContent() {
	regDto := suite.RegisterForTest("a", "a")

	responseStatusCode, _ := suite.SendGet(
		"/api/user/orders",
		regDto.Token,
	)

	assert.Equal(suite.T(), http.StatusNoContent, responseStatusCode)
}

func (suite *ordersTestSuite) TestOrders401IfNoToken() {
	suite.RegisterForTest("a", "a")

	responseStatusCode, _ := suite.SendGet(
		"/api/user/orders",
		"",
	)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseStatusCode)
}

func (suite *ordersTestSuite) Test200IfAlreadyUploaded() {
	var err error
	regDto := suite.RegisterForTest("a", "a")
	_, err = suite.Repository.AddOrder(
		context.Background(),
		"12345678903",
		regDto.ID,
		"st1", nil,
		time.Time{},
	)
	assert.NoError(suite.T(), err)

	responseStatusCode := suite.SendPost(
		"/api/user/orders",
		"text/plain",
		regDto.Token,
		bytes.NewBuffer([]byte("12345678903")),
	)

	require.Equal(suite.T(), http.StatusOK, responseStatusCode)
}

func (suite *ordersTestSuite) Test409IfAlreadyUploadedByAnotherUser() {
	var err error
	anotherUser := suite.RegisterForTest("another", "a")

	_, err = suite.Repository.AddOrder(
		context.Background(),
		"12345678903",
		anotherUser.ID,
		"st1", nil,
		time.Time{},
	)
	assert.NoError(suite.T(), err)

	regDto := suite.RegisterForTest("a", "a")
	responseStatusCode := suite.SendPost(
		"/api/user/orders",
		"text/plain",
		regDto.Token,
		bytes.NewBuffer([]byte("12345678903")),
	)

	require.Equal(suite.T(), http.StatusConflict, responseStatusCode)
}

func (suite *ordersTestSuite) Test422IfInvalidNumber() {
	var _ error
	regDto := suite.RegisterForTest("a", "a")

	responseStatusCode := suite.SendPost(

		"/api/user/orders",
		"text/plain",
		regDto.Token,
		bytes.NewBuffer([]byte("1234567890")),
	)

	require.Equal(suite.T(), http.StatusUnprocessableEntity, responseStatusCode)
}
