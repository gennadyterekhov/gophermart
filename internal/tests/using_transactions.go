package tests

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/storage"
	"github.com/stretchr/testify/assert"
)

const TestDBDSN = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"

type (
	beforeOrAfterFunc            func(*testing.T, *storage.DB)
	testCase                     func(*testing.T)
	TestRunnerWithBeforeAndAfter func(testCase) testCase
)

type SuiteUsingTransactions struct {
	db *storage.DB
}

func NewSuiteUsingTransactions(db *storage.DB) *SuiteUsingTransactions {
	return &SuiteUsingTransactions{
		db: db,
	}
}

func (suite *SuiteUsingTransactions) SetDB(db *storage.DB) {
	suite.db = db
}

func (suite *SuiteUsingTransactions) CustomBeforeEach() TestRunnerWithBeforeAndAfter {
	return setBeforeAndAfterEach(beforeEach, afterEach, suite.db)
}

// UsingTransactions rollbacks all db changes after each test
func (suite *SuiteUsingTransactions) UsingTransactions() TestRunnerWithBeforeAndAfter {
	return setBeforeAndAfterEach(beforeEach, afterEach, suite.db)
}

// UsingTransactions rollbacks all db changes after each test
func UsingTransactions() TestRunnerWithBeforeAndAfter {
	// TODO fix to use existing?
	return setBeforeAndAfterEach(beforeEach, afterEach, storage.NewDB(TestDBDSN))
}

// setBeforeAndAfterEach takes before and after functions and returns a function called by t.Run().
func setBeforeAndAfterEach(beforeFunc, afterFunc beforeOrAfterFunc, db *storage.DB) TestRunnerWithBeforeAndAfter {
	return func(test testCase) testCase {
		return func(t *testing.T) {
			if beforeFunc != nil {
				beforeFunc(t, db)
			}

			test(t)

			if afterFunc != nil {
				afterFunc(t, db)
			}
		}
	}
}

func beforeEach(t *testing.T, db *storage.DB) {
	if db == nil {
		panic("db is nil in using_transactions.beforeEach")
	}
	var err error

	db.Connection.Tx, err = db.Connection.BeginTx(context.Background(), nil)
	assert.NoError(t, err)

	_, err = db.Connection.Exec("SAVEPOINT test")
	assert.NoError(t, err)
}

func afterEach(t *testing.T, db *storage.DB) {
	if db == nil {
		panic("db is nil in using_transactions.beforeEach")
	}
	var err error

	_, err = db.Connection.Exec("ROLLBACK TO SAVEPOINT test")
	assert.NoError(t, err)
	err = db.Connection.Rollback()
	assert.NoError(t, err)

	assert.NoError(t, err)
}
