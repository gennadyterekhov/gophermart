package withdrawals

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"

	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/stretchr/testify/assert"

	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"

	"github.com/gennadyterekhov/gophermart/internal/tests"
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
