package withdrawals

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/balance"
	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests.BeforeAll()
	code := m.Run()
	tests.AfterAll()
	os.Exit(code)
}

func TestCanGetWithdrawals(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		userDto := helpers.RegisterForTest("a", "a")
		withdrawalNewest, withdrawalMedium, withdrawalOldest := createDifferentWithdrawals(t, userDto)

		ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
		all, err := GetAll(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(*all))
		assert.Equal(t, withdrawalOldest.ID, (*all)[0].ID)
		assert.Equal(t, withdrawalMedium.ID, (*all)[1].ID)
		assert.Equal(t, withdrawalNewest.ID, (*all)[2].ID)
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

func TestCanCreateWithdrawals(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		userDto := helpers.RegisterForTest("a", "a")
		var accrual int64 = 101
		_, err := repositories.AddOrder(
			context.Background(),
			"a",
			userDto.ID,
			"",
			&accrual,
			time.Time{},
		)
		require.NoError(t, err)

		ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
		reqDto := &requests.Withdrawals{
			Order: "a",
			Sum:   1,
		}
		_, err = Create(ctx, reqDto)
		assert.NoError(t, err)

		bal, _ := balance.GetBalance(context.Background(), userDto.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), bal)
	}))
}

func TestCannotCreateWithdrawalsIfNotEnoughBalance(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		userDto := helpers.RegisterForTest("a", "a")
		var accrual int64 = 5
		_, err := repositories.AddOrder(
			context.Background(),
			"a",
			userDto.ID,
			"",
			&accrual,
			time.Time{},
		)
		require.NoError(t, err)

		ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
		reqDto := &requests.Withdrawals{
			Order: "a",
			Sum:   10,
		}
		_, err = Create(ctx, reqDto)
		assert.Equal(t, ErrorInsufficientFunds, err.Error())
	}))
}

func createDifferentWithdrawals(
	t *testing.T,
	userDto *responses.Register,
) (*models.Withdrawal, *models.Withdrawal, *models.Withdrawal) {
	withdrawalNewest, err := repositories.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 0,
		time.Time{},
	)
	assert.NoError(t, err)
	withdrawalMedium, err := repositories.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 0,
		time.Time{}.AddDate(-1, 0, 0),
	)
	assert.NoError(t, err)
	withdrawalOldest, err := repositories.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 0,
		time.Time{}.AddDate(-10, 0, 0),
	)
	assert.NoError(t, err)
	return withdrawalNewest, withdrawalMedium, withdrawalOldest
}
