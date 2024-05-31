package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanSendOrdersRequest(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		var err error
		regDto := helpers.RegisterForTest("a", "a")
		_, err = repositories.AddOrder(
			context.Background(),
			"1",
			regDto.ID,
			"st1", nil,
			time.Time{},
		)
		assert.NoError(t, err)

		var tenDollarsOrThousandCents int64 = 1000
		_, err = repositories.AddOrder(
			context.Background(),
			"2",
			regDto.ID,
			"st2", &tenDollarsOrThousandCents,
			time.Time{}.AddDate(1, 0, 0),
		)
		assert.NoError(t, err)

		responseStatusCode, bodyAsBytes := tests.SendGet(
			t,
			tests.TestServer,
			"/api/user/orders",
			regDto.Token,
		)

		require.Equal(t, http.StatusOK, responseStatusCode)
		responseBody := make([]models.OrderFloats, 0)

		err = json.Unmarshal(bodyAsBytes, &responseBody)
		assert.NoError(t, err)
		require.Equal(t, 2, len(responseBody))
		assert.Nil(t, (responseBody)[0].Accrual)
		assert.Equal(t, float64(10), *(responseBody)[1].Accrual)
	}))
}

func TestOrders204IfNoContent(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		regDto := helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := tests.SendGet(
			t,
			tests.TestServer,
			"/api/user/orders",
			regDto.Token,
		)

		assert.Equal(t, http.StatusNoContent, responseStatusCode)
	}))
}

func TestOrders401IfNoToken(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := tests.SendGet(
			t,
			tests.TestServer,
			"/api/user/orders",
			"",
		)

		assert.Equal(t, http.StatusUnauthorized, responseStatusCode)
	}))
}
