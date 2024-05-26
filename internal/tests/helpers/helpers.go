package helpers

import (
	"context"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
)

func RegisterForTest(login string, password string) *responses.Register {
	reqDto := &requests.Register{Login: login, Password: password}
	resDto, err := register.Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}
