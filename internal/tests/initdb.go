package tests

import (
	"context"
	"database/sql"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/config"
)

func InitDB() (*sql.DB, *sql.Tx) {
	const testDBDsn = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"
	config.ServerConfig.DBDsn = testDBDsn
	storage.Connection = storage.CreateDBStorage(testDBDsn)

	transaction, err := storage.Connection.DBConnection.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return storage.Connection.DBConnection, transaction
}
