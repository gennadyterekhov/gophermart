package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanSendWithdrawalsRequest(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		withdrawalNewest, err := repositories.AddWithdrawal(
			context.Background(),
			regDto.ID,
			"", 0,
			time.Time{},
		)
		assert.NoError(t, err)

		responseStatusCode, bodyAsBytes := tests.SendGet(
			t,
			tests.TestServer,
			"/api/user/withdrawals",
			regDto.Token,
		)

		require.Equal(t, http.StatusOK, responseStatusCode)
		responseBody := &responses.Withdrawals{}
		err = json.Unmarshal(bodyAsBytes, responseBody)
		assert.NoError(t, err)
		require.Equal(t, 1, len(*responseBody))
		assert.Equal(t, withdrawalNewest.ID, (*responseBody)[0].ID)
	}))
}

func Test204IfNoContent(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := tests.SendGet(
			t,
			tests.TestServer,
			"/api/user/withdrawals",
			regDto.Token,
		)

		assert.Equal(t, http.StatusNoContent, responseStatusCode)
	}))
}

func TestCannotCreateWithdrawalsWithIncorrectNumber(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		var accrual int64 = 10
		_, err := repositories.AddOrder(
			context.Background(),
			"a",
			regDto.ID,
			"",
			&accrual,
			time.Time{},
		)
		require.NoError(t, err)

		rawJSON := `{"order":"123", "sum":1}`
		responseStatusCode := tests.SendPost(
			t,
			tests.TestServer,
			"/api/user/balance/withdraw",
			"application/json",
			regDto.Token,
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusUnprocessableEntity, responseStatusCode)
	}))
}

func Test402WhenNotEnoughBalance(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")

		rawJSON := `{"order":"4417123456789113", "sum":1}`
		responseStatusCode := tests.SendPost(
			t,
			tests.TestServer,
			"/api/user/balance/withdraw",
			"application/json",
			regDto.Token,
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusPaymentRequired, responseStatusCode)
	}))
}

func Test401IfNoToken(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := tests.SendGet(
			t,
			tests.TestServer,
			"/api/user/withdrawals",
			"",
		)

		assert.Equal(t, http.StatusUnauthorized, responseStatusCode)
	}))
}
