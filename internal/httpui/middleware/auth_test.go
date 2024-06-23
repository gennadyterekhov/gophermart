package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/gennadyterekhov/gophermart/internal/tests"
)

type authTestSuite struct {
	suite.Suite
	tests.TestHTTPServer
	tests.SuiteUsingTransactions
}

func TestAuth(t *testing.T) {
	server := httptest.NewServer(
		getTestRouter(),
	)
	suiteInstance := &authTestSuite{}
	suiteInstance.SetDB(storage.NewDB(helpers.TestDBDSN))
	suiteInstance.Server = server
	suite.Run(t, suiteInstance)
}

func getTestRouter() *chi.Mux {
	testRouter := chi.NewRouter()
	testRouter.Get(
		"/auth",
		WithAuth(
			http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
			}),
		).ServeHTTP,
	)
	testRouter.Post(
		"/json",
		RequestContentTypeJSON(
			http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
			}),
		).ServeHTTP,
	)
	testRouter.Post(
		"/luhn",
		Luhn(
			http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
			}),
		).ServeHTTP,
	) //
	return testRouter
}

func (suite *authTestSuite) TestCanAuthWithToken() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		resDto := helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := suite.SendGet(
			"/auth",
			resDto.Token,
		)

		assert.Equal(t, http.StatusOK, responseStatusCode)
	}))
}

func (suite *authTestSuite) Test401IfNoToken() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		responseStatusCode, _ := suite.SendGet(
			"/auth",
			"incorrect token",
		)

		assert.Equal(t, http.StatusUnauthorized, responseStatusCode)
	}))
}
