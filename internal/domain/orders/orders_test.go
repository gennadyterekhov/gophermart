package orders

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/base"

	"github.com/gennadyterekhov/gophermart/internal/client"
	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	base.BaseSuite
	Service Service
}

func newSuite() *testSuite {
	suiteInstance := &testSuite{}
	base.InitBaseSuite(suiteInstance)
	suiteInstance.Service = NewService(suiteInstance.GetRepository(), client.NewClient("", suiteInstance.GetRepository()))

	return suiteInstance
}

func TestDomainOrders(t *testing.T) {
	suite.Run(t, newSuite())
}

func (suite *testSuite) TestCanGetOrders() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		userDto := suite.RegisterForTest("a", "a")
		withdrawalNewest, withdrawalMedium, withdrawalOldest := suite.createDifferentOrders(userDto)

		ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
		all, err := suite.Service.GetAll(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(*all))
		assert.Equal(t, withdrawalOldest.Number, (*all)[0].Number)
		assert.Equal(t, withdrawalMedium.Number, (*all)[1].Number)
		assert.Equal(t, withdrawalNewest.Number, (*all)[2].Number)
	}))
}

func (suite *testSuite) TestNoContentReturnsError() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		userDto := suite.RegisterForTest("a", "a")
		ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
		_, err := suite.Service.GetAll(ctx)
		assert.Error(t, err)
		assert.Equal(suite.T(), err.Error(), ErrorNoContent)
	}))
}

func (suite *testSuite) TestCanCreateOrder() {
}

func (suite *testSuite) TestCanOrderStatusIsAutomaticallyUpdated() {
}

func (suite *testSuite) createDifferentOrders(
	userDto *responses.Register,
) (*order.Order, *order.Order, *order.Order) {
	var ten int64 = 10
	withdrawalNewest, err := suite.Service.Repository.AddOrder(
		context.Background(),
		"1",
		userDto.ID,
		"", &ten,
		time.Time{},
	)
	assert.NoError(suite.T(), err)
	withdrawalMedium, err := suite.Service.Repository.AddOrder(
		context.Background(),
		"2",
		userDto.ID,
		"", &ten,
		time.Time{}.AddDate(-1, 0, 0),
	)
	assert.NoError(suite.T(), err)
	withdrawalOldest, err := suite.Service.Repository.AddOrder(
		context.Background(),
		"3",
		userDto.ID,
		"", &ten,
		time.Time{}.AddDate(-10, 0, 0),
	)
	assert.NoError(suite.T(), err)
	return withdrawalNewest, withdrawalMedium, withdrawalOldest
}
