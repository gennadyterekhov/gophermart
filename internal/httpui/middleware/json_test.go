package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

type jsonTestSuite struct {
	suite.Suite
	tests.TestHTTPServer
}

func TestJSON(t *testing.T) {
	server := httptest.NewServer(
		getTestRouter(),
	)
	suiteInstance := &jsonTestSuite{}
	suiteInstance.Server = server
	suite.Run(t, suiteInstance)
}

func (suite *jsonTestSuite) TestCanSendIfJson() {
	suite.T().Run("", func(t *testing.T) {
		path := "/json"
		req, err := http.NewRequest(http.MethodPost, suite.Server.URL+path, nil)
		require.NoError(suite.T(), err)
		req.Header.Set("Content-Type", "application/json")

		response, err := suite.Server.Client().Do(req)
		require.NoError(suite.T(), err)
		response.Body.Close()

		assert.Equal(suite.T(), http.StatusOK, response.StatusCode)
	})
}

func (suite *jsonTestSuite) Test400IfNotJson() {
	suite.T().Run("", func(t *testing.T) {
		path := "/json"
		req, err := http.NewRequest(http.MethodPost, suite.Server.URL+path, nil)
		require.NoError(suite.T(), err)

		response, err := suite.Server.Client().Do(req)
		require.NoError(suite.T(), err)
		response.Body.Close()

		assert.Equal(suite.T(), http.StatusBadRequest, response.StatusCode)
	})
}
