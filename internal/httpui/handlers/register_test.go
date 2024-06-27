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

type testSuite struct {
	with_server.BaseSuiteWithServer
}

func TestRegister(t *testing.T) {
	suiteInstance := &testSuite{}
	with_server.InitBaseSuiteWithServer(suiteInstance)
	suite.Run(t, suiteInstance)
}

func (suite *testSuite) TestCanSendRegisterRequest() {
	cases := []struct {
		name        string
		contentType string

		code int
	}{
		{
			name:        "ok",
			contentType: "application/json",
			code:        http.StatusOK,
		},
		{
			name:        "400",
			contentType: "application",
			code:        http.StatusBadRequest,
		},
	}

	for _, tt := range cases {
		rawJSON := `{"login":"a", "password":"a"}`
		responseStatusCode, bodyAsBytes := suite.SendPostAndReturnBody(
			"/api/user/register",
			tt.contentType,
			"",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(suite.T(), tt.code, responseStatusCode)
		if tt.code == http.StatusOK {
			responseBody := &responses.Register{ID: 0, Token: ""}
			err := json.Unmarshal(bodyAsBytes, responseBody)
			assert.NoError(suite.T(), err)
			assert.NotEqual(suite.T(), "", responseBody.Token)
			assert.NotEqual(suite.T(), 0, responseBody.Token)
		}
	}
}

func (suite *testSuite) Test409IfSameLogin() {
	suite.RegisterForTest("a", "a")

	rawJSON := `{"login":"a", "password":"b"}`
	responseStatusCode := suite.SendPost(
		"/api/user/register",
		"application/json",
		"",
		bytes.NewBuffer([]byte(rawJSON)),
	)

	assert.Equal(suite.T(), http.StatusConflict, responseStatusCode)
}
