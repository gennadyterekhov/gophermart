package tests

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/base"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type orderRepositoryTest struct {
	base.BaseSuite
}

func (suite *orderRepositoryTest) SetupSuite() {
	base.InitBaseSuite(suite)
}

func (suite *orderRepositoryTest) TestCanGetOrdersFromRepo() {
	suite.T().Run("", func(t *testing.T) {
		regDto := suite.RegisterForTest("a", "a")
		orderNewest, orderMedium, orderOldest := createDifferentOrders(suite.T(), suite.Repository, regDto)

		orders, err := suite.Repository.GetAllOrdersForUser(context.Background(), regDto.ID)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 3, len(orders))
		assert.Equal(suite.T(), orderOldest.Number, orders[0].Number)
		assert.Equal(suite.T(), orderMedium.Number, orders[1].Number)
		assert.Equal(suite.T(), orderNewest.Number, orders[2].Number)
	})
}

func (suite *orderRepositoryTest) TestCanInsertOrder() {
	suite.T().Run("", func(t *testing.T) {
		regDto := suite.RegisterForTest("a", "a")
		_, err := suite.Repository.AddOrder(context.Background(), "1", regDto.ID, "", nil, time.Time{})
		assert.NoError(suite.T(), err)
	})
}

func TestOrdersRepo(t *testing.T) {
	suite.Run(t, new(orderRepositoryTest))
}

func createDifferentOrders(
	t *testing.T,
	repo *repositories.RepositoryMock,
	userDto *responses.Register,
) (*order.Order, *order.Order, *order.Order) {
	orderNewest, err := repo.AddOrder(
		context.Background(),
		"1",
		userDto.ID,
		"", nil,
		time.Time{},
	)
	assert.NoError(t, err)
	orderMedium, err := repo.AddOrder(
		context.Background(),
		"2",
		userDto.ID,
		"", nil,
		time.Time{}.AddDate(-1, 0, 0),
	)
	assert.NoError(t, err)
	orderOldest, err := repo.AddOrder(
		context.Background(),
		"3",
		userDto.ID,
		"", nil,
		time.Time{}.AddDate(-10, 0, 0),
	)
	assert.NoError(t, err)
	return orderNewest, orderMedium, orderOldest
}
