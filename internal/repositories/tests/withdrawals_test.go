package tests

import (
	"context"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/base"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type withdrawalsRepositoryTest struct {
	base.BaseSuite
}

func newSuite() *withdrawalsRepositoryTest {
	suiteInstance := &withdrawalsRepositoryTest{}
	base.InitBaseSuite(suiteInstance)

	return suiteInstance
}

func (suite *withdrawalsRepositoryTest) TestCanInsertAndGetAllWithdrawals() {
	suite.T().Run("", func(t *testing.T) {
		var err error
		user, err := suite.Repository.AddUser(context.Background(), "a", "a")
		assert.NoError(t, err)

		user2, err := suite.Repository.AddUser(context.Background(), "b", "a")
		assert.NoError(t, err)

		_, err = suite.Repository.AddWithdrawal(context.Background(), user.ID, "a", 0, time.Time{})
		assert.NoError(t, err)
		_, err = suite.Repository.AddWithdrawal(context.Background(), user.ID, "b", 0, time.Time{})
		assert.NoError(t, err)
		_, err = suite.Repository.AddWithdrawal(context.Background(), user2.ID, "b", 0, time.Time{})
		assert.NoError(t, err)

		wdrs, _ := suite.Repository.GetAllWithdrawalsForUser(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(wdrs))

		wdrs2, _ := suite.Repository.GetAllWithdrawalsForUser(context.Background(), user2.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(wdrs2))
	})
}

func TestWithdrawals(t *testing.T) {
	suite.Run(t, newSuite())
}
