package tests

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestCanInsertAndGetAllWithdrawals(t *testing.T) {
	run := tests.UsingTransactions()
	t.Run("", run(func(t *testing.T) {
		var err error
		user, err := repositories.AddUser(context.Background(), "a", "a")
		assert.NoError(t, err)

		user2, err := repositories.AddUser(context.Background(), "b", "a")
		assert.NoError(t, err)

		_, err = repositories.AddWithdrawal(context.Background(), user.ID, "a", 0, time.Time{})
		assert.NoError(t, err)
		_, err = repositories.AddWithdrawal(context.Background(), user.ID, "b", 0, time.Time{})
		assert.NoError(t, err)
		_, err = repositories.AddWithdrawal(context.Background(), user2.ID, "b", 0, time.Time{})
		assert.NoError(t, err)

		wdrs, _ := repositories.GetAllWithdrawalsForUser(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(wdrs))

		wdrs2, _ := repositories.GetAllWithdrawalsForUser(context.Background(), user2.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(wdrs2))
	}))
}
