package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/server"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/stretchr/testify/assert"
)

type loginTestSuite struct {
	server.BaseSuiteWithServer
}

func TestLoginHandler(t *testing.T) {
	suiteInstance := &loginTestSuite{}
	server.InitBaseSuiteWithServer(suiteInstance)

	suite.Run(t, suiteInstance)
}

func (suite *loginTestSuite) TestCanSendLoginRequest() {
	suite.RegisterForTest("a", "a")

	rawJSON := `{"login":"a", "password":"a"}`
	responseStatusCode, bodyAsBytes := suite.SendPostAndReturnBody(
		"/api/user/login",
		"application/json",
		"",
		bytes.NewBuffer([]byte(rawJSON)),
	)

	assert.Equal(suite.T(), http.StatusOK, responseStatusCode)

	responseBody := &responses.Login{Token: ""}
	err := json.Unmarshal(bodyAsBytes, responseBody)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), "", responseBody.Token)
}

func (suite *loginTestSuite) TestCannotLoginWithWrongFieldName() {
	suite.RegisterForTest("a", "a")

	rawJSON := `{"logi":"a", "password":"a"}`
	responseStatusCode := suite.SendPost(
		"/api/user/login",
		"application/json",
		"",
		bytes.NewBuffer([]byte(rawJSON)),
	)

	assert.Equal(suite.T(), http.StatusBadRequest, responseStatusCode)
}

func (suite *loginTestSuite) TestCannotLoginWithWrongContentType() {
	suite.RegisterForTest("a", "a")

	rawJSON := `{"login":"a", "password":"a"}`
	responseStatusCode := suite.SendPost(
		"/api/user/login",
		"application",
		"",
		bytes.NewBuffer([]byte(rawJSON)),
	)

	assert.Equal(suite.T(), http.StatusBadRequest, responseStatusCode)
}
