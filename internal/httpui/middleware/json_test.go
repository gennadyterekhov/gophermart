package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/gennadyterekhov/gophermart/internal/tests"
)

func TestCanSendIfJson(t *testing.T) {
	run := tests.UsingTransactions()
	setupTestServer()

	t.Run("", run(func(t *testing.T) {
		path := "/json"
		req, err := http.NewRequest(http.MethodPost, tests.TestServer.URL+path, nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		response, err := tests.TestServer.Client().Do(req)
		require.NoError(t, err)
		response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode)
	}))
}

func Test400IfNotJson(t *testing.T) {
	run := tests.UsingTransactions()
	setupTestServer()

	t.Run("", run(func(t *testing.T) {
		path := "/json"
		req, err := http.NewRequest(http.MethodPost, tests.TestServer.URL+path, nil)
		require.NoError(t, err)

		response, err := tests.TestServer.Client().Do(req)
		require.NoError(t, err)
		response.Body.Close()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	}))
}
