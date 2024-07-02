package storage

import (
	"context"
	"database/sql"

	"github.com/gennadyterekhov/gophermart/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type QueryMaker interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Ping() error
	Close() error
}

type DB struct {
	Connection QueryMaker
}

func NewDB(dsn string) *DB {
	logger.CustomLogger.Debugln("opening database connection with dsn ", dsn)

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.CustomLogger.Debugln("could not connect to db using dsn: " + dsn + " " + err.Error())
		panic(err)
	}

	return &DB{
		Connection: conn,
	}
}
