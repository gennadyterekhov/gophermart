package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tests.BeforeAll()
	code := m.Run()
	tests.AfterAll()
	os.Exit(code)
}

func TestCatGetUserById(t *testing.T) {
	run := tests.UsingTransactions()
	t.Run("", run(func(t *testing.T) {
		user, err := AddUser(context.Background(), "a", "a")
		assert.NoError(t, err)

		user, err = GetUserById(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "a", user.Login)
		assert.Equal(t, "a", user.Password)
	}))
}

func TestCatInsertUser(t *testing.T) {
	run := tests.UsingTransactions()
	t.Run("", run(func(t *testing.T) {
		user, err := AddUser(context.Background(), "a", "a")
		assert.NoError(t, err)
		assert.Equal(t, "a", user.Login)
		assert.Equal(t, "a", user.Password)
	}))
}
