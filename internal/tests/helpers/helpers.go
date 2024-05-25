package helpers

import (
	"context"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
)

func RegisterForTest(login string, password string) *responses.Register {
	reqDto := &requests.Register{Login: login, Password: password}
	resDto, err := auth.Register(context.Background(), reqDto)
	if err != nil {
		panic(err)
	}
	return resDto
}
