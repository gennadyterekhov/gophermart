package middleware

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gennadyterekhov/gophermart/internal/tests"
)

func TestLuhnOk(t *testing.T) {
	run := tests.UsingTransactions()
	setupTestServer()

	t.Run("", run(func(t *testing.T) {
		rawJSON := `{"order":"4417123456789113"}`

		responseStatusCode := tests.SendPostWithoutToken(
			t,
			tests.TestServer,
			"/luhn",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusOK, responseStatusCode)
	}))
}

func TestLuhnOkWhenTextPlain(t *testing.T) {
	run := tests.UsingTransactions()
	setupTestServer()

	t.Run("", run(func(t *testing.T) {
		responseStatusCode := tests.SendPost(
			t,
			tests.TestServer,
			"/luhn",
			"text/plain",
			"",
			bytes.NewBuffer([]byte("12345678903")),
		)

		assert.Equal(t, http.StatusOK, responseStatusCode)
	}))
}

func Test422WhenNoOrderInBody(t *testing.T) {
	run := tests.UsingTransactions()
	setupTestServer()

	t.Run("", run(func(t *testing.T) {
		rawJSON := `{"hello":"4417123456789113"}`

		responseStatusCode := tests.SendPostWithoutToken(
			t,
			tests.TestServer,
			"/luhn",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusUnprocessableEntity, responseStatusCode)
	}))
}

func Test422WhenIncorrectNumber(t *testing.T) {
	run := tests.UsingTransactions()
	setupTestServer()

	t.Run("", run(func(t *testing.T) {
		rawJSON := `{"order":"4417123456789119"}`

		responseStatusCode := tests.SendPostWithoutToken(
			t,
			tests.TestServer,
			"/luhn",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusUnprocessableEntity, responseStatusCode)
	}))
	t.Run("", run(func(t *testing.T) {
		rawJSON := `{"order":"441712a456789113"}`

		responseStatusCode := tests.SendPostWithoutToken(
			t,
			tests.TestServer,
			"/luhn",
			bytes.NewBuffer([]byte(rawJSON)),
		)

		assert.Equal(t, http.StatusUnprocessableEntity, responseStatusCode)
	}))
}
