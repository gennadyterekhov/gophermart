package balance

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/base"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/stretchr/testify/assert"
)

type testSuite struct {
	base.BaseSuite
	Service Service
}

func (suite *testSuite) SetupSuite() {
	base.InitBaseSuite(suite)
	suite.Service = NewService(suite.GetRepository())
}

func Test(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) TestCanGetBalance() {
	userDto := suite.RegisterForTest("a", "a")
	suite.createDifferentWithdrawals(userDto)
	var startBalance int64 = 10
	_, err := suite.Service.Repository.AddOrder(context.Background(), "", userDto.ID, "", &startBalance, time.Time{})
	assert.NoError(suite.T(), err)

	ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
	resDto, err := suite.Service.GetBalanceResponse(ctx)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), int64(10-(1+2+3)), resDto.Current) // TODO fix currency float
	assert.Equal(suite.T(), int64(1+2+3), resDto.Withdrawn)
}

func (suite *testSuite) createDifferentWithdrawals(
	userDto *responses.Register,
) (*models.Withdrawal, *models.Withdrawal, *models.Withdrawal) {
	withdrawalNewest, err := suite.Service.Repository.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 1,
		time.Time{},
	)
	assert.NoError(suite.T(), err)
	withdrawalMedium, err := suite.Service.Repository.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 2,
		time.Time{}.AddDate(-1, 0, 0),
	)
	assert.NoError(suite.T(), err)
	withdrawalOldest, err := suite.Service.Repository.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 3,
		time.Time{}.AddDate(-10, 0, 0),
	)
	assert.NoError(suite.T(), err)
	return withdrawalNewest, withdrawalMedium, withdrawalOldest
}
