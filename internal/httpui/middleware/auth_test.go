package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/gennadyterekhov/gophermart/internal/tests"
)

func TestMain(m *testing.M) {
	tests.BeforeAll()
	code := m.Run()
	tests.AfterAll()
	os.Exit(code)
}

func setupTestServer() {
	testRouter := chi.NewRouter()
	testRouter.Get(
		"/test",
		WithAuth(
			http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
			}),
		).ServeHTTP,
	)

	tests.TestServer = httptest.NewServer(
		testRouter,
	)
}

func TestCanAuthWithToken(t *testing.T) {
	run := tests.UsingTransactions()
	setupTestServer()

	t.Run("", run(func(t *testing.T) {
		resDto := helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := tests.SendGet(
			t,
			tests.TestServer,
			"/test",
			resDto.Token,
		)

		assert.Equal(t, http.StatusOK, responseStatusCode)
	}))
}

func Test401IfNoToken(t *testing.T) {
	run := tests.UsingTransactions()
	setupTestServer()

	t.Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := tests.SendGet(
			t,
			tests.TestServer,
			"/test",
			"incorrect token",
		)

		assert.Equal(t, http.StatusUnauthorized, responseStatusCode)
	}))
}
