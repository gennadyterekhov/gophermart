package auth

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/auth/token"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/assert"
)

type loginTest struct {
	suite.Suite
	tests.SuiteUsingTransactions
	Service Service
}

func (suite *loginTest) SetupSuite() {
	db := storage.NewDB(helpers.TestDBDSN)
	repo := repositories.NewRepository(db)
	suite.SetDB(db)
	//	suiteInstance.SetDB(storage.NewDB(helpers.TestDBDSN))
	suite.Service = NewService(repo)
}

func TestLogin(t *testing.T) {
	suite.Run(t, new(loginTest))
}

func (suite *loginTest) TestCanLogin() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		registerForTest("a", "a")

		reqDto := &requests.Login{Login: "a", Password: "a"}
		resDto, err := suite.Service.Login(context.Background(), reqDto)
		assert.NoError(t, err)

		err = token.ValidateToken(resDto.Token, "a")
		assert.NoError(t, err)
	}))
}

func (suite *loginTest) TestCannotLoginWithWrongLogin() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		registerForTest("a", "a")

		reqDto := &requests.Login{Login: "b", Password: "a"}
		_, err := suite.Service.Login(context.Background(), reqDto)
		assert.Error(t, err)
		assert.Equal(t, ErrorWrongCredentials, err.Error())
	}))
}

func (suite *loginTest) TestCannotLoginWithWrongPassword() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		registerForTest("a", "a")

		reqDto := &requests.Login{Login: "a", Password: "b"}
		_, err := suite.Service.Login(context.Background(), reqDto)
		assert.Error(t, err)
		assert.Equal(t, ErrorWrongCredentials, err.Error())
	}))
}

func registerForTest(login string, password string) *responses.Register {
	db := storage.NewDB(tests.TestDBDSN)
	repo := repositories.NewRepository(db)
	registerService := register.NewService(repo)
	reqDto := &requests.Register{Login: login, Password: password}
	resDto, err := registerService.Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}
