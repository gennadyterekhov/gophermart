package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type luhnTestSuite struct {
	suite.Suite
	tests.TestHTTPServer
}

func TestLuhn(t *testing.T) {
	server := httptest.NewServer(
		getTestRouter(),
	)
	suiteInstance := &luhnTestSuite{}

	suiteInstance.Server = server
	suite.Run(t, suiteInstance)
}

func (suite *luhnTestSuite) TestLuhnOk() {
	suite.T().Run("", func(t *testing.T) {
		rawJSON := `{"order":"4417123456789113"}`

		responseStatusCode := suite.SendPostWithoutToken(
			"/luhn",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(suite.T(), http.StatusOK, responseStatusCode)
	})
}

func (suite *luhnTestSuite) TestLuhnOkWhenTextPlain() {
	suite.T().Run("", func(t *testing.T) {
		responseStatusCode := suite.SendPost(
			"/luhn",
			"text/plain",
			"",
			bytes.NewBuffer([]byte("12345678903")),
		)

		assert.Equal(suite.T(), http.StatusOK, responseStatusCode)
	})
}

func (suite *luhnTestSuite) Test422WhenNoOrderInBody() {
	suite.T().Run("", func(t *testing.T) {
		rawJSON := `{"hello":"4417123456789113"}`

		responseStatusCode := suite.SendPostWithoutToken(
			"/luhn",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(suite.T(), http.StatusUnprocessableEntity, responseStatusCode)
	})
}

func (suite *luhnTestSuite) Test422WhenIncorrectNumber() {
	suite.T().Run("", func(t *testing.T) {
		rawJSON := `{"order":"4417123456789119"}`

		responseStatusCode := suite.SendPostWithoutToken(
			"/luhn",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(suite.T(), http.StatusUnprocessableEntity, responseStatusCode)
	})
	suite.T().Run("", func(t *testing.T) {
		rawJSON := `{"order":"441712a456789113"}`

		responseStatusCode := suite.SendPostWithoutToken(
			"/luhn",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(suite.T(), http.StatusUnprocessableEntity, responseStatusCode)
	})
}
