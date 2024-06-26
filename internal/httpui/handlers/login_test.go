package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/with_server"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/stretchr/testify/assert"
)

type loginTestSuite struct {
	with_server.BaseSuiteWithServer
}

func TestLoginHandler(t *testing.T) {
	suiteInstance := &loginTestSuite{}
	with_server.InitBaseSuiteWithServer(suiteInstance)

	suite.Run(t, suiteInstance)
}

func (suite *loginTestSuite) TestCanSendLoginRequest() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		suite.RegisterForTest("a", "a")

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
		suite.RegisterForTest("a", "a")

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
		suite.RegisterForTest("a", "a")

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
