package storage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

type dbTest struct {
	suite.Suite
}

func initDB() *DB {
	const TestDBDSN = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"
	return NewDB(TestDBDSN)
}

func (suite *dbTest) TestDbExists() {
	_db := initDB()
	err := _db.Connection.Ping()
	assert.NoError(suite.T(), err)
}

func (suite *dbTest) TestDbTableExists() {
	_db := initDB()
	rawSQLString := "select * from users limit 1;"
	_, err := _db.Connection.Exec(rawSQLString)
	assert.NoError(suite.T(), err)

	rawSQLString = "select * from orders limit 1;"
	_, err = _db.Connection.Exec(rawSQLString)
	assert.NoError(suite.T(), err)

	rawSQLString = "select * from withdrawals limit 1;"
	_, err = _db.Connection.Exec(rawSQLString)
	assert.NoError(suite.T(), err)
}

func TestDb(t *testing.T) {
	suite.Run(t, new(dbTest))
}
