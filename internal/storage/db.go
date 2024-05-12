package storage

import (
	"database/sql"

	"github.com/gennadyterekhov/gophermart/internal/storage/migration"

	"github.com/gennadyterekhov/gophermart/internal/config"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

type DBStorage struct {
	DBConnection *sql.DB
}

var Connection = CreateDBStorage()

func CreateDBStorage() *DBStorage {
	conn, err := sql.Open("pgx", config.ServerConfig.DBDsn)
	if err != nil {
		logger.ZapSugarLogger.Panicln("could not connect to db using dsn: " + config.ServerConfig.DBDsn)
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
