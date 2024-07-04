package withdrawals

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/base"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/domain/balance"
	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSuite struct {
	base.BaseSuite
	Service Service
}

func (suite *testSuite) SetupSuite() {
	base.InitBaseSuite(suite)

	suite.Service = NewService(suite.GetRepository(), balance.NewService(suite.GetRepository()))
}

func TestWithdrawalsDomain(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) TestCanGetWithdrawals() {
	userDto := suite.RegisterForTest("a", "a")
	withdrawalNewest, withdrawalMedium, withdrawalOldest := suite.createDifferentWithdrawals(userDto)

	ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
	all, err := suite.Service.GetAll(ctx)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 3, len(*all))
	assert.Equal(suite.T(), withdrawalOldest.ID, (*all)[0].ID)
	assert.Equal(suite.T(), withdrawalMedium.ID, (*all)[1].ID)
	assert.Equal(suite.T(), withdrawalNewest.ID, (*all)[2].ID)
}

func (suite *testSuite) TestNoContentReturnsError() {
	userDto := suite.RegisterForTest("a", "a")
	ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
	_, err := suite.Service.GetAll(ctx)
	assert.Equal(suite.T(), err.Error(), ErrorNoContent)
}

func (suite *testSuite) TestCanCreateWithdrawals() {
	userDto := suite.RegisterForTest("a", "a")
	var accrual int64 = 101
	_, err := suite.Service.Repository.AddOrder(
		context.Background(),
		"a",
		userDto.ID,
		"",
		&accrual,
		time.Time{},
	)
	require.NoError(suite.T(), err)

	ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
	reqDto := &requests.Withdrawals{
		Order: "a",
		Sum:   1,
	}
	_, err = suite.Service.Create(ctx, reqDto)
	assert.NoError(suite.T(), err)

	bal, _ := suite.Service.BalanceService.GetBalance(context.Background(), userDto.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), bal)
}

func (suite *testSuite) TestCannotCreateWithdrawalsIfNotEnoughBalance() {
	userDto := suite.RegisterForTest("a", "a")
	var accrual int64 = 5
	_, err := suite.Service.Repository.AddOrder(
		context.Background(),
		"a",
		userDto.ID,
		"",
		&accrual,
		time.Time{},
	)
	require.NoError(suite.T(), err)

	ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, userDto.ID)
	reqDto := &requests.Withdrawals{
		Order: "a",
		Sum:   10,
	}
	_, err = suite.Service.Create(ctx, reqDto)
	assert.Equal(suite.T(), ErrorInsufficientFunds, err.Error())
}

func (suite *testSuite) createDifferentWithdrawals(
	userDto *responses.Register,
) (*models.Withdrawal, *models.Withdrawal, *models.Withdrawal) {
	withdrawalNewest, err := suite.Service.Repository.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 0,
		time.Time{},
	)
	assert.NoError(suite.T(), err)
	withdrawalMedium, err := suite.Service.Repository.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 0,
		time.Time{}.AddDate(-1, 0, 0),
	)
	assert.NoError(suite.T(), err)
	withdrawalOldest, err := suite.Service.Repository.AddWithdrawal(
		context.Background(),
		userDto.ID,
		"", 0,
		time.Time{}.AddDate(-10, 0, 0),
	)
	assert.NoError(suite.T(), err)
	return withdrawalNewest, withdrawalMedium, withdrawalOldest
}
