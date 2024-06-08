package tests

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/config"
	"github.com/gennadyterekhov/gophermart/internal/storage"
	"github.com/stretchr/testify/assert"
)

func InitDB() storage.QueryMaker {
	const testDBDsn = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"
	config.ServerConfig.DBDsn = testDBDsn
	storage.DBClient = storage.CreateDBStorage(testDBDsn)
	return storage.DBClient.Connection
}

func BeforeAll() {
	logger.Init()
	_ = InitDB()
	storage.DBClient.Connection.UseTx = true
}

func AfterAll() {
	storage.DBClient.Connection.UseTx = false
	storage.DBClient.Connection.Close()
}

func beforeEach(t *testing.T) {
	var err error
	storage.DBClient.Connection.Tx, err = storage.DBClient.Connection.BeginTx(context.Background(), nil)
	assert.NoError(t, err)

	_, err = storage.DBClient.Connection.Exec("SAVEPOINT test")
	assert.NoError(t, err)
}

func afterEach(t *testing.T) {
	var err error

	_, err = storage.DBClient.Connection.Exec("ROLLBACK TO SAVEPOINT test")
	assert.NoError(t, err)
	err = storage.DBClient.Connection.Rollback()
	assert.NoError(t, err)

	assert.NoError(t, err)
}

// setBeforeAndAfterEach takes before and after functions and returns a function called by t.Run().
func setBeforeAndAfterEach(beforeFunc, afterFunc func(*testing.T)) func(func(*testing.T)) func(*testing.T) {
	return func(test func(*testing.T)) func(*testing.T) {
		return func(t *testing.T) {
			if beforeFunc != nil {
				beforeFunc(t)
			}
			test(t)
			if afterFunc != nil {
				afterFunc(t)
			}
		}
	}
}

// UsingTransactions rollbacks all db changes after each test
func UsingTransactions() func(func(*testing.T)) func(*testing.T) {
	return setBeforeAndAfterEach(beforeEach, afterEach)
}
