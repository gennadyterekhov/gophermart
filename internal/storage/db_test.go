package storage

import (
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func initDB() {
	const testDBDsn = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"
	config.ServerConfig.DBDsn = testDBDsn
	DBClient = CreateDBStorage(testDBDsn)
}

func TestDbExists(t *testing.T) {
	initDB()
	err := DBClient.Connection.Ping()
	assert.NoError(t, err)
}

func TestDbTableExists(t *testing.T) {
	initDB()
	rawSQLString := "select * from users limit 1;"
	_, err := DBClient.Connection.Exec(rawSQLString)
	assert.NoError(t, err)

	rawSQLString = "select * from orders limit 1;"
	_, err = DBClient.Connection.Exec(rawSQLString)
	assert.NoError(t, err)

	rawSQLString = "select * from withdrawals limit 1;"
	_, err = DBClient.Connection.Exec(rawSQLString)
	assert.NoError(t, err)
}
