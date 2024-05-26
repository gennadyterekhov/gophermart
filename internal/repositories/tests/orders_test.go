package tests

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCanGetOrders(t *testing.T) {
	run := tests.UsingTransactions()
	t.Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		orderNewest, orderMedium, orderOldest := createDifferentOrders(t, regDto)

		orders, err := repositories.GetAllOrdersForUser(context.Background(), regDto.ID)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(orders))
		assert.Equal(t, orderOldest.Number, orders[0].Number)
		assert.Equal(t, orderMedium.Number, orders[1].Number)
		assert.Equal(t, orderNewest.Number, orders[2].Number)
	}))
}

func TestCanInsertOrder(t *testing.T) {
	run := tests.UsingTransactions()
	t.Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		_, err := repositories.AddOrder(context.Background(), "1", regDto.ID, "", nil, time.Time{})
		assert.NoError(t, err)
	}))
}

func createDifferentOrders(
	t *testing.T,
	userDto *responses.Register,
) (*models.Order, *models.Order, *models.Order) {
	orderNewest, err := repositories.AddOrder(
		context.Background(),
		"1",
		userDto.ID,
		"", nil,
		time.Time{},
	)
	assert.NoError(t, err)
	orderMedium, err := repositories.AddOrder(
		context.Background(),
		"2",
		userDto.ID,
		"", nil,
		time.Time{}.AddDate(-1, 0, 0),
	)
	assert.NoError(t, err)
	orderOldest, err := repositories.AddOrder(
		context.Background(),
		"3",
		userDto.ID,
		"", nil,
		time.Time{}.AddDate(-10, 0, 0),
	)
	assert.NoError(t, err)
	return orderNewest, orderMedium, orderOldest
}
