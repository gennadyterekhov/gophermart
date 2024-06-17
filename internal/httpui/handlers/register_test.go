package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/config"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/register"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

type testSuite struct {
	suite.Suite
	tests.SuiteUsingTransactions
	tests.TestHTTPServer
}

func Test(t *testing.T) {
	suiteInstance := &testSuite{}
	suiteInstance.SetDB(storage.NewDB(helpers.TestDBDSN))
	suite.Run(t, suiteInstance)
}

func (suite *testSuite) TestCanSendRegisterRequest() {
	run := suite.UsingTransactions()
	db := storage.NewDB(tests.TestDBDSN)
	tests.InitTestServer(NewRouter(config.NewConfig(), db).Router)
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
		suite.T().Run(tt.name, run(func(t *testing.T) {
			rawJSON := `{"login":"a", "password":"a"}`
			responseStatusCode, bodyAsBytes := suite.SendPostAndReturnBody(
				"/api/user/register",
				tt.contentType,
				"",
				bytes.NewBuffer([]byte(rawJSON)),
			)

			assert.Equal(t, tt.code, responseStatusCode)
			if tt.code == http.StatusOK {
				responseBody := &responses.Register{ID: 0, Token: ""}
				err := json.Unmarshal(bodyAsBytes, responseBody)
				assert.NoError(t, err)
				assert.NotEqual(t, "", responseBody.Token)
				assert.NotEqual(t, 0, responseBody.Token)
			}
		}))
	}
}

func (suite *testSuite) Test409IfSameLogin() {
	run := suite.UsingTransactions()
	db := storage.NewDB(tests.TestDBDSN)
	tests.InitTestServer(NewRouter(config.NewConfig(), db).Router)
	repo := repositories.NewRepository(db)
	service := register.NewService(repo)
	suite.T().Run("", run(func(t *testing.T) {
		reqDto := &requests.Register{
			Login:    "a",
			Password: "a",
		}
		_, err := service.Register(context.Background(), reqDto)
		assert.NoError(t, err)

		rawJSON := `{"login":"a", "password":"b"}`
		responseStatusCode := suite.SendPost(
			"/api/user/register",
			"application/json",
			"",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusConflict, responseStatusCode)
	}))
}
