package tests

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

type userRepositoryTest struct {
	suite.Suite
	tests.SuiteUsingTransactions
	Repository repositories.Repository
}

func (suite *userRepositoryTest) SetupSuite() {
	db := storage.NewDB(helpers.TestDBDSN)
	suite.SetDB(db)
	suite.Repository = repositories.NewRepository(db)
}

func (suite *userRepositoryTest) TestCanGetUserByID() {
	run := suite.UsingTransactions()
	repo := repositories.NewRepository(storage.NewDB(tests.TestDBDSN))

	suite.T().Run("", run(func(t *testing.T) {
		user, err := repo.AddUser(context.Background(), "a", "a")
		assert.NoError(t, err)

		user, err = repo.GetUserByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "a", user.Login)
		assert.Equal(t, "a", user.Password)
	}))
}

func (suite *userRepositoryTest) TestCanInsertUser() {
	run := suite.UsingTransactions()
	repo := repositories.NewRepository(storage.NewDB(tests.TestDBDSN))
	// TODO https://github.com/gennadyterekhov/gophermart/issues/53
	suite.T().Run("", run(func(t *testing.T) {
		user, err := repo.AddUser(context.Background(), "a", "a")
		assert.NoError(t, err)
		assert.Equal(t, "a", user.Login)
		assert.Equal(t, "a", user.Password)
	}))
}

func TestUser(t *testing.T) {
	suite.Run(t, new(userRepositoryTest))
}
