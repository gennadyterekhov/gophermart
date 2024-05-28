package handlers

import (
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

func TestCanSendBalanceRequest(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")
		_, err := repositories.AddWithdrawal(context.Background(), regDto.ID, "", 100, time.Time{})
		assert.NoError(t, err)

		var startBalance int64 = 1000 // cents => $10
		_, err = repositories.AddOrder(context.Background(), "", regDto.ID, "", &startBalance, time.Time{})
		assert.NoError(t, err)

		responseStatusCode, bodyAsBytes := tests.SendGet(
			t,
			tests.TestServer,
			"/api/user/balance",
			regDto.Token,
		)

		require.Equal(t, http.StatusOK, responseStatusCode)
		responseBody := &responses.BalanceExternal{}
		err = json.Unmarshal(bodyAsBytes, responseBody)
		assert.NoError(t, err)
		require.Equal(t, float64(1), responseBody.Withdrawn)
		assert.Equal(t, float64(9), responseBody.Current)
	}))
}

func TestBalance401IfNoToken(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := tests.SendGet(
			t,
			tests.TestServer,
			"/api/user/balance",
			"",
		)

		assert.Equal(t, http.StatusUnauthorized, responseStatusCode)
	}))
}
