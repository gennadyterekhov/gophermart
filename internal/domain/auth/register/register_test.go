package register

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/stretchr/testify/assert"
)

type testSuite struct {
	suite.Suite
	tests.SuiteUsingTransactions
	Service Service
}

func (suite *testSuite) SetupSuite() {
	const TestDBDSN = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"

	db := storage.NewDB(TestDBDSN)
	repo := repositories.NewRepository(db)
	suite.SetDB(db)
	suite.Service = NewService(repo)
}

func Test(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) TestCanRegister() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		reqDto := &requests.Register{
			Login:    "a",
			Password: "a",
		}
		resDto, err := suite.Service.Register(context.Background(), reqDto)
		assert.NoError(t, err)
		assert.NotEqual(t, "", resDto.Token)

		user, err := suite.Service.Repository.GetUserByID(context.Background(), resDto.ID)
		assert.NoError(t, err)
		assert.Equal(t, "a", user.Login)
		assert.NotEqual(t, "a", user.Password)
	}))
}

func (suite *testSuite) TestCannotRegisterWhenLoginAlreadyUsed() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		var err error
		_, err = suite.Service.Register(context.Background(), &requests.Register{Login: "a", Password: "a"})
		assert.NoError(t, err)
		_, err = suite.Service.Register(context.Background(), &requests.Register{Login: "a", Password: "a"})
		assert.Equal(t, "ERROR: duplicate key value violates unique constraint \"users_login_key\" (SQLSTATE 23505)", err.Error())
	}))
}
