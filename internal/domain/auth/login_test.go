package auth

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/base"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/token"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/stretchr/testify/assert"
)

type loginTest struct {
	base.BaseSuite
	Service Service
}

func newSuite() *loginTest {
	suiteInstance := &loginTest{}
	base.InitBaseSuite(suiteInstance)
	suiteInstance.Service = NewService(suiteInstance.GetRepository())

	return suiteInstance
}

func TestLogin(t *testing.T) {
	suite.Run(t, newSuite())
}

func (suite *loginTest) TestCanLogin() {
	suite.RegisterForTest("a", "a")

	reqDto := &requests.Login{Login: "a", Password: "a"}
	resDto, err := suite.Service.Login(context.Background(), reqDto)
	assert.NoError(suite.T(), err)

	err = token.ValidateToken(resDto.Token, "a")
	assert.NoError(suite.T(), err)
}

func (suite *loginTest) TestCannotLoginWithWrongLogin() {
	suite.RegisterForTest("a", "a")

	reqDto := &requests.Login{Login: "b", Password: "a"}
	_, err := suite.Service.Login(context.Background(), reqDto)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrorWrongCredentials, err.Error())
}

func (suite *loginTest) TestCannotLoginWithWrongPassword() {
	suite.RegisterForTest("a", "a")

	reqDto := &requests.Login{Login: "a", Password: "b"}
	_, err := suite.Service.Login(context.Background(), reqDto)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrorWrongCredentials, err.Error())
}
