package auth

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestCanLogin(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		userDto := registerForTest("a", "a")

		reqDto := &requests.Login{Login: "a", Password: "a"}
		resDto, err := Login(context.Background(), reqDto)
		assert.NoError(t, err)

		err = validateToken(resDto.Token, userDto.ID)
		assert.NoError(t, err)
	}))
}

func TestCannotLoginWithWrongLogin(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		registerForTest("a", "a")

		reqDto := &requests.Login{Login: "b", Password: "a"}
		_, err := Login(context.Background(), reqDto)
		assert.Error(t, err)
		assert.Equal(t, ErrorWrongCredentials, err.Error())
	}))
}

func TestCannotLoginWithWrongPassword(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		registerForTest("a", "a")

		reqDto := &requests.Login{Login: "a", Password: "b"}
		_, err := Login(context.Background(), reqDto)
		assert.Error(t, err)
		assert.Equal(t, ErrorWrongCredentials, err.Error())
	}))
}

func registerForTest(login string, password string) *responses.Register {
	reqDto := &requests.Register{Login: login, Password: password}
	resDto, err := Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}
