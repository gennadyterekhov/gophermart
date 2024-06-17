package tests

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/assert"
)

type orderRepositoryTest struct {
	suite.Suite
	tests.SuiteUsingTransactions
	Repository repositories.Repository
}

func (suite *orderRepositoryTest) SetupSuite() {
	db := storage.NewDB(helpers.TestDBDSN)
	suite.SetDB(db)
	suite.Repository = repositories.NewRepository(db)
}

func (suite *orderRepositoryTest) TestCanGetOrders() {
	run := suite.UsingTransactions()
	suite.T().Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		orderNewest, orderMedium, orderOldest := createDifferentOrders(t, &suite.Repository, regDto)

		orders, err := suite.Repository.GetAllOrdersForUser(context.Background(), regDto.ID)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(orders))
		assert.Equal(t, orderOldest.Number, orders[0].Number)
		assert.Equal(t, orderMedium.Number, orders[1].Number)
		assert.Equal(t, orderNewest.Number, orders[2].Number)
	}))
}

func (suite *orderRepositoryTest) TestCanInsertOrder() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		_, err := suite.Repository.AddOrder(context.Background(), "1", regDto.ID, "", nil, time.Time{})
		assert.NoError(t, err)
	}))
}

func Test(t *testing.T) {
	suite.Run(t, new(orderRepositoryTest))
}

func createDifferentOrders(
	t *testing.T,
	repo *repositories.Repository,
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
