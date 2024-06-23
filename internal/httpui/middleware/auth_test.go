package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/gennadyterekhov/gophermart/internal/tests"

	"github.com/stretchr/testify/suite"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type authTestSuite struct {
	suite.Suite
	tests.TestHTTPServer
	Repository *repositories.RepositoryMock
}

func (suite *authTestSuite) SetupSuite() {
	suite.Repository = repositories.NewRepositoryMock()
	suite.Server = httptest.NewServer(
		getTestRouter(),
	)
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(authTestSuite))
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
	resDto := suite.registerForTest("a", "a")

	responseStatusCode, _ := suite.SendGet(
		"/auth",
		resDto.Token,
	)

	assert.Equal(suite.T(), http.StatusOK, responseStatusCode)
}

func (suite *authTestSuite) Test401IfNoToken() {
	suite.registerForTest("a", "a")

	responseStatusCode, _ := suite.SendGet(
		"/auth",
		"incorrect token",
	)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseStatusCode)
}

func (suite *authTestSuite) registerForTest(login string, password string) *responses.Register {
	reqDto := &requests.Register{Login: login, Password: password}
	service := register.NewService(suite.Repository)
	resDto, err := service.Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}
