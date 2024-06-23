package tests

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/storage"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

type withdrawalsRepositoryTest struct {
	suite.Suite
	tests.SuiteUsingTransactions
	Repository repositories.Repository
}

func (suite *withdrawalsRepositoryTest) SetupSuite() {
	db := storage.NewDB(helpers.TestDBDSN)
	suite.SetDB(db)
	suite.Repository = repositories.NewRepository(db)
}

func (suite *withdrawalsRepositoryTest) TestCanInsertAndGetAllWithdrawals() {
	run := suite.UsingTransactions()
	repo := repositories.NewRepository(storage.NewDB(tests.TestDBDSN))

	suite.T().Run("", run(func(t *testing.T) {
		var err error
		user, err := repo.AddUser(context.Background(), "a", "a")
		assert.NoError(t, err)

		user2, err := repo.AddUser(context.Background(), "b", "a")
		assert.NoError(t, err)

		_, err = repo.AddWithdrawal(context.Background(), user.ID, "a", 0, time.Time{})
		assert.NoError(t, err)
		_, err = repo.AddWithdrawal(context.Background(), user.ID, "b", 0, time.Time{})
		assert.NoError(t, err)
		_, err = repo.AddWithdrawal(context.Background(), user2.ID, "b", 0, time.Time{})
		assert.NoError(t, err)

		wdrs, _ := repo.GetAllWithdrawalsForUser(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(wdrs))

		wdrs2, _ := repo.GetAllWithdrawalsForUser(context.Background(), user2.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(wdrs2))
	}))
}

func TestWithdrawals(t *testing.T) {
	suite.Run(t, new(withdrawalsRepositoryTest))
}
