package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

// needed to avoid import cycle
func initDB() (*sql.DB, *sql.Tx) {
	const testDBDsn = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"
	config.ServerConfig.DBDsn = testDBDsn
	Connection = CreateDBStorage(testDBDsn)

	transaction, err := Connection.DBConnection.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return Connection.DBConnection, transaction
}

func TestDbExists(t *testing.T) {
	dbConnection, tx := initDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	err := dbConnection.Ping()
	assert.NoError(t, err)
}

func TestDbTableExists(t *testing.T) {
	dbConnection, tx := initDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	rawSQLString := "select * from users limit 1;"
	_, err := tx.Exec(rawSQLString)
	assert.NoError(t, err)

	rawSQLString = "select * from orders limit 1;"
	_, err = tx.Exec(rawSQLString)
	assert.NoError(t, err)

	rawSQLString = "select * from withdrawals limit 1;"
	_, err = tx.Exec(rawSQLString)
	assert.NoError(t, err)
}
