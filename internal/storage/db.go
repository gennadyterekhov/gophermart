package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gennadyterekhov/gophermart/internal/config"

	"github.com/gennadyterekhov/gophermart/internal/storage/migration"

	"github.com/gennadyterekhov/gophermart/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type QueryMaker interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Ping() error
	Close() error
	Commit() error
	Rollback() error
}

func (ct *ConnectionOrTransaction) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if ct.UseTx {
		return ct.Tx.QueryContext(ctx, query, args...)
	}
	return ct.Conn.QueryContext(ctx, query, args...)
}

func (ct *ConnectionOrTransaction) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if ct.UseTx {
		return ct.Tx.QueryRowContext(ctx, query, args...)
	}
	return ct.Conn.QueryRowContext(ctx, query, args...)
}

func (ct *ConnectionOrTransaction) Ping() error {
	return ct.Conn.Ping()
}

func (ct *ConnectionOrTransaction) Close() error {
	return ct.Conn.Close()
}

func (ct *ConnectionOrTransaction) Exec(query string, args ...any) (sql.Result, error) {
	if ct.UseTx {
		return ct.Tx.Exec(query, args...)
	}
	return ct.Conn.Exec(query, args...)
}

func (ct *ConnectionOrTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if ct.UseTx {
		return ct.Tx.Exec(query, args...)
	}
	return ct.Conn.ExecContext(ctx, query, args...)
}

func (ct *ConnectionOrTransaction) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if ct.UseTx {
		if ct.Tx == nil {
			var err error
			ct.Tx, err = ct.Conn.BeginTx(ctx, opts)
			return ct.Tx, err
		}
		return nil, fmt.Errorf("beginning transaction from existing transaction")
	}
	return ct.Conn.BeginTx(ctx, opts)
}

func (ct *ConnectionOrTransaction) Commit() error {
	if ct.UseTx {
		err := ct.Tx.Commit()
		ct.Tx = nil
		return err
	}
	return nil
}

func (ct *ConnectionOrTransaction) Rollback() error {
	if ct.UseTx {
		err := ct.Tx.Rollback()
		ct.Tx = nil
		return err
	}
	return nil
}

type ConnectionOrTransaction struct {
	Conn  *sql.DB
	Tx    *sql.Tx
	UseTx bool
}

type DB struct {
	Connection *ConnectionOrTransaction
}

var DBClient = createDefaultDBClient()

func createDefaultDBClient() *DB {
	return CreateDBStorage(config.ServerConfig.DBDsn)
}

func CreateDBStorage(dsn string) *DB {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.CustomLogger.Debugln("could not connect to db using dsn: " + dsn + " " + err.Error())
		panic(err)
	}

	migration.RunMigrations(conn)

	ct := &ConnectionOrTransaction{
		Conn:  conn,
		Tx:    nil,
		UseTx: false,
	}

	return &DB{
		Connection: ct,
	}
}
