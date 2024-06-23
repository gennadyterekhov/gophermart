package helpers

import (
	"context"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/storage"
	"github.com/gennadyterekhov/gophermart/internal/tests"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
)

const TestDBDSN = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"

func RegisterForTest(login string, password string) *responses.Register {
	reqDto := &requests.Register{Login: login, Password: password}
	service := register.NewService(repositories.NewRepository(storage.NewDB(tests.TestDBDSN)))
	resDto, err := service.Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}
