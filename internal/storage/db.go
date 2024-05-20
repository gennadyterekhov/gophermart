package storage

import (
	"database/sql"

	"github.com/gennadyterekhov/gophermart/internal/config"

	"github.com/gennadyterekhov/gophermart/internal/storage/migration"

	"github.com/gennadyterekhov/gophermart/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	DBConnection *sql.DB
}

var Connection = CreateDefaultDBStorage()

func CreateDefaultDBStorage() *DBStorage {
	return CreateDBStorage(config.ServerConfig.DBDsn)
}

func CreateDBStorage(dsn string) *DBStorage {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.ZapSugarLogger.Panicln("could not connect to db using dsn: " + dsn + " " + err.Error())
	}

	migration.RunMigrations(conn)

	return &DBStorage{
		DBConnection: conn,
	}
}

func (strg *DBStorage) CloseDB() error {
	err := strg.DBConnection.Close()
	if err != nil {
		logger.ZapSugarLogger.Errorln("could not close db", err.Error())
	}
	return err
}

func (strg *DBStorage) GetDB() *DBStorage {
	return strg
}
