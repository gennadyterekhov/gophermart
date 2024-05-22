package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tests.BeforeAll()
	code := m.Run()
	tests.AfterAll()
	os.Exit(code)
}

func TestCanSendRegisterRequest(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

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
		t.Run(tt.name, run(func(t *testing.T) {
			rawJSON := `{"login":"a", "password":"a"}`
			responseStatusCode, bodyAsBytes := tests.SendPostAndReturnBody(
				t,
				tests.TestServer,
				"/api/user/register",
				tt.contentType,
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

func Test409IfSameLogin(t *testing.T) {
	run := tests.UsingTransactions()
	tests.InitTestServer(GetRouter())

	t.Run("", run(func(t *testing.T) {
		reqDto := &requests.Register{
			Login:    "a",
			Password: "a",
		}
		_, err := auth.Register(context.Background(), reqDto)
		assert.NoError(t, err)

		rawJSON := `{"login":"a", "password":"b"}`
		responseStatusCode := tests.SendPost(
			t,
			tests.TestServer,
			"/api/user/register",
			"application/json",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusConflict, responseStatusCode)
	}))
}
