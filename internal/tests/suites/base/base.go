package base

import (
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests/traits"
	"github.com/stretchr/testify/suite"
)

type HasRepo interface {
	SetRepository(repo *repositories.RepositoryMock)
	GetRepository() *repositories.RepositoryMock
}

type HasLifecycleMethods interface {
	SetupTest()
	TearDownTest()
}

type CanRegister interface {
	RegisterForTest(login string, password string) *responses.Register
}

type UsingTransactions interface {
	UsingTransactions() traits.TestRunnerWithBeforeAndAfter
	HasLifecycleMethods
}

type BaseSuiteInterface interface {
	UsingTransactions
	HasLifecycleMethods
	CanRegister
	HasRepo
}

type BaseSuite struct {
	suite.Suite
	traits.CanRegisterAndUsingTransactions
}

func InitBaseSuite[T BaseSuiteInterface](srv T) {
	repo := repositories.NewRepositoryMock()
	srv.SetRepository(repo)
}

func (s *BaseSuite) SetRepository(repo *repositories.RepositoryMock) {
	s.Repository = repo
}

func (s *BaseSuite) GetRepository() *repositories.RepositoryMock {
	return s.Repository
}
