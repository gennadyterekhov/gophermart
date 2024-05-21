package auth

import (
	"context"
	"os"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tests.BeforeAll()
	code := m.Run()
	tests.AfterAll()
	os.Exit(code)
}

func TestCanRegister(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		reqDto := &requests.Register{
			Login:    "a",
			Password: "a",
		}
		resDto, err := Register(context.Background(), reqDto)
		assert.NoError(t, err)
		assert.NotEqual(t, "", resDto.Token)

		user, err := repositories.GetUserById(context.Background(), resDto.ID)
		assert.NoError(t, err)
		assert.Equal(t, "a", user.Login)
		assert.NotEqual(t, "a", user.Password)
	}))
}

func TestCannotRegisterWhenLoginAlreadyUsed(t *testing.T) {
	run := tests.UsingTransactions()

	t.Run("", run(func(t *testing.T) {
		var err error
		_, err = Register(context.Background(), &requests.Register{Login: "a", Password: "a"})
		assert.NoError(t, err)
		_, err = Register(context.Background(), &requests.Register{Login: "a", Password: "a"})
		assert.Equal(t, "ERROR: duplicate key value violates unique constraint \"users_login_key\" (SQLSTATE 23505)", err.Error())
	}))
}
