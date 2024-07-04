package traits

import (
	"context"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

type CanRegister struct {
	Repository *repositories.RepositoryMock
}

func (suite *CanRegister) SetupTest() {
	suite.Repository.Clear()
}

func (suite *CanRegister) TearDownTest() {
	suite.Repository.Clear()
}

func (suite *CanRegister) RegisterForTest(login string, password string) *responses.Register {
	reqDto := &requests.Register{Login: login, Password: password}
	service := register.NewService(suite.Repository)
	resDto, err := service.Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}
