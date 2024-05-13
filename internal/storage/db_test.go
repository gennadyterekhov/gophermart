package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

const TestDBDsn = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"

func initDB() (*sql.DB, *sql.Tx) {
	config.ServerConfig.DBDsn = TestDBDsn
	dbConnection, err := sql.Open("pgx", TestDBDsn)
	if err != nil {
		panic(err)
	}

	transaction, err := dbConnection.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return dbConnection, transaction
}

func TestDbExists(t *testing.T) {
	t.Skip("only manual use because depends on host")

	dbConnection, tx := initDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	err := dbConnection.Ping()
	assert.NoError(t, err)
}

func TestDbTableExists(t *testing.T) {
	t.Skip("only manual use because depends on host")

	dbConnection, tx := initDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	rawSQLString := "select * from users limit 1;"
	_, err := tx.Exec(rawSQLString)
	assert.NoError(t, err)

	rawSQLString = "select * from balances limit 1;"
	_, err = tx.Exec(rawSQLString)
	assert.NoError(t, err)

	rawSQLString = "select * from orders limit 1;"
	_, err = tx.Exec(rawSQLString)
	assert.NoError(t, err)

	rawSQLString = "select * from withdrawals limit 1;"
	_, err = tx.Exec(rawSQLString)
	assert.NoError(t, err)
}
