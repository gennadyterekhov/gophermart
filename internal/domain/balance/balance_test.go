package balance

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
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

func TestCanGetBalance(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		userDto := helpers.RegisterForTest("a", "a")
		createDifferentWithdrawals(t, userDto)
		var startBalance int64 = 10
		_, err := repositories.AddOrder(context.Background(), "", userDto.ID, "", &startBalance, time.Time{})
		assert.NoError(t, err)

		ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
		resDto, err := GetBalanceResponse(ctx)
		assert.NoError(t, err)

		assert.Equal(t, int64(10-(1+2+3)), resDto.Current) // TODO fix currency float
		assert.Equal(t, int64(1+2+3), resDto.Withdrawn)
	}))
}

func createDifferentWithdrawals(
	t *testing.T,
	userDto *responses.Register,
) (*models.Withdrawal, *models.Withdrawal, *models.Withdrawal) {
	withdrawalNewest, err := repositories.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 1,
		time.Time{},
	)
	assert.NoError(t, err)
	withdrawalMedium, err := repositories.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 2,
		time.Time{}.AddDate(-1, 0, 0),
	)
	assert.NoError(t, err)
	withdrawalOldest, err := repositories.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 3,
		time.Time{}.AddDate(-10, 0, 0),
	)
	assert.NoError(t, err)
	return withdrawalNewest, withdrawalMedium, withdrawalOldest
}
