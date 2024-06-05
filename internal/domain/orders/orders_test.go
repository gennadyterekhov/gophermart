package orders

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tests.BeforeAll()
	code := m.Run()
	tests.AfterAll()
	os.Exit(code)
}

func TestCanGetOrders(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		userDto := helpers.RegisterForTest("a", "a")
		withdrawalNewest, withdrawalMedium, withdrawalOldest := createDifferentOrders(t, userDto)

		ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
		all, err := GetAll(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(*all))
		assert.Equal(t, withdrawalOldest.Number, (*all)[0].Number)
		assert.Equal(t, withdrawalMedium.Number, (*all)[1].Number)
		assert.Equal(t, withdrawalNewest.Number, (*all)[2].Number)
	}))
}

func TestNoContentReturnsError(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		userDto := helpers.RegisterForTest("a", "a")
		ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
		_, err := GetAll(ctx)
		assert.Equal(t, err.Error(), ErrorNoContent)
	}))
}

func createDifferentOrders(
	t *testing.T,
	userDto *responses.Register,
) (*order.Order, *order.Order, *order.Order) {
	var ten int64 = 10
	withdrawalNewest, err := repositories.AddOrder(
		context.Background(),
		"1",
		userDto.ID,
		"", &ten,
		time.Time{},
	)
	assert.NoError(t, err)
	withdrawalMedium, err := repositories.AddOrder(
		context.Background(),
		"2",
		userDto.ID,
		"", &ten,
		time.Time{}.AddDate(-1, 0, 0),
	)
	assert.NoError(t, err)
	withdrawalOldest, err := repositories.AddOrder(
		context.Background(),
		"3",
		userDto.ID,
		"", &ten,
		time.Time{}.AddDate(-10, 0, 0),
	)
	assert.NoError(t, err)
	return withdrawalNewest, withdrawalMedium, withdrawalOldest
}
