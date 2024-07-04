package register

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/stretchr/testify/assert"
)

type testSuite struct {
	// cannot use 	base.BaseSuite because of import cycle
	suite.Suite
	Repository *repositories.RepositoryMock
	Service    Service
}

func (suite *testSuite) SetupSuite() {
	suite.Repository = repositories.NewRepositoryMock()
	suite.Service = NewService(suite.Repository)
}

func (suite *testSuite) SetupTest() {
	suite.Repository.Clear()
}

func (suite *testSuite) TearDownTest() {
	suite.Repository.Clear()
}

func TestRegisterDomain(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) TestCanRegister() {
	reqDto := &requests.Register{
		Login:    "a",
		Password: "a",
	}
	resDto, err := suite.Service.Register(context.Background(), reqDto)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), "", resDto.Token)

	user, err := suite.Service.Repository.GetUserByID(context.Background(), resDto.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "a", user.Login)
	assert.NotEqual(suite.T(), "a", user.Password)
}

func (suite *testSuite) TestCannotRegisterWhenLoginAlreadyUsed() {
	var err error
	_, err = suite.Service.Register(context.Background(), &requests.Register{Login: "a", Password: "a"})
	assert.NoError(suite.T(), err)
	_, err = suite.Service.Register(context.Background(), &requests.Register{Login: "a", Password: "a"})
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "ERROR: duplicate key value violates unique constraint \"users_login_key\" (SQLSTATE 23505)", err.Error())
}
