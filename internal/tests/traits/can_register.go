package traits

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/storage"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

type (
	beforeOrAfterFunc            func(*testing.T, *storage.DB)
	testCase                     func(*testing.T)
	TestRunnerWithBeforeAndAfter func(testCase) testCase
)

type CanRegisterAndUsingTransactions struct {
	Repository *repositories.RepositoryMock
}

// deprecated
func (suite *CanRegisterAndUsingTransactions) UsingTransactions() TestRunnerWithBeforeAndAfter {
	return setBeforeAndAfterEach(beforeEach, afterEach, nil)
}

func (suite *CanRegisterAndUsingTransactions) SetupTest() {
	suite.Repository.Clear()
}

func (suite *CanRegisterAndUsingTransactions) TearDownTest() {
	suite.Repository.Clear()
}

func (suite *CanRegisterAndUsingTransactions) RegisterForTest(login string, password string) *responses.Register {
	reqDto := &requests.Register{Login: login, Password: password}
	service := register.NewService(suite.Repository)
	resDto, err := service.Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}

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
}

func afterEach(t *testing.T, db *storage.DB) {
}
