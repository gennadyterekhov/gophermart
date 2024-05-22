package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestCanSendLoginRequest(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		registerForTest("a", "a")

		rawJSON := `{"login":"a", "password":"a"}`
		responseStatusCode, bodyAsBytes := tests.SendPostAndReturnBody(
			t,
			tests.TestServer,
			"/api/user/login",
			"application/json",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusOK, responseStatusCode)

		responseBody := &responses.Login{Token: ""}
		err := json.Unmarshal(bodyAsBytes, responseBody)
		assert.NoError(t, err)
		assert.NotEqual(t, "", responseBody.Token)
	}))
}

func TestCannotLoginWithWrongFieldName(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		registerForTest("a", "a")

		rawJSON := `{"logi":"a", "password":"a"}`
		responseStatusCode := tests.SendPost(
			t,
			tests.TestServer,
			"/api/user/login",
			"application/json",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusBadRequest, responseStatusCode)
	}))
}

func TestCannotLoginWithWrongContentType(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		registerForTest("a", "a")

		rawJSON := `{"login":"a", "password":"a"}`
		responseStatusCode := tests.SendPost(
			t,
			tests.TestServer,
			"/api/user/login",
			"application",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusBadRequest, responseStatusCode)
	}))
}

func registerForTest(login string, password string) *responses.Register {
	reqDto := &requests.Register{Login: login, Password: password}
	resDto, err := auth.Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}
