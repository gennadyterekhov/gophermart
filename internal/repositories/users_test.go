package repositories

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestCatGetUserById(t *testing.T) {
	dbConnection, tx := tests.InitDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	rawSQLString := "insert into users (id, login, password) values (1, 'a', 'a');"
	_, err := tx.Exec(rawSQLString)
	assert.NoError(t, err)

	user, err := GetUserByIdTx(context.Background(), tx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "a", user.Login)
	assert.Equal(t, "a", user.Password)
}

func TestCatInsertUser(t *testing.T) {
	dbConnection, tx := tests.InitDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	user, err := AddUserTx(context.Background(), tx, "a", "a")
	assert.NoError(t, err)
	assert.Equal(t, "a", user.Login)
	assert.Equal(t, "a", user.Password)
}
