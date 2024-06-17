package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

type loginTestSuite struct {
	suite.Suite
	tests.SuiteUsingTransactions
	tests.TestHTTPServer
	Repository repositories.Repository
}

func TestLogin(t *testing.T) {
	db := storage.NewDB(tests.TestDBDSN)

	suiteInstance := &loginTestSuite{
		Repository: repositories.NewRepository(db),
	}
	suiteInstance.SetDB(db)

	suite.Run(t, suiteInstance)
}

func (suite *loginTestSuite) TestCanSendLoginRequest() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		rawJSON := `{"login":"a", "password":"a"}`
		responseStatusCode, bodyAsBytes := suite.SendPostAndReturnBody(
			"/api/user/login",
			"application/json",
			"",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusOK, responseStatusCode)

		responseBody := &responses.Login{Token: ""}
		err := json.Unmarshal(bodyAsBytes, responseBody)
		assert.NoError(t, err)
		assert.NotEqual(t, "", responseBody.Token)
	}))
}

func (suite *loginTestSuite) TestCannotLoginWithWrongFieldName() {
	run := suite.UsingTransactions()
	suite.T().Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		rawJSON := `{"logi":"a", "password":"a"}`
		responseStatusCode := suite.SendPost(
			"/api/user/login",
			"application/json",
			"",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusBadRequest, responseStatusCode)
	}))
}

func (suite *loginTestSuite) TestCannotLoginWithWrongContentType() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		helpers.RegisterForTest("a", "a")

		rawJSON := `{"login":"a", "password":"a"}`
		responseStatusCode := suite.SendPost(
			"/api/user/login",
			"application",
			"",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusBadRequest, responseStatusCode)
	}))
}
