package auth

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/logger"
	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/alexedwards/argon2id"
)

const ErrorWrongCredentials = "unknown credentials"

func Login(ctx context.Context, reqDto *requests.Login) (*responses.Login, error) {
	user, err := repositories.GetUserByLogin(ctx, reqDto.Login)
	if err != nil {
		return nil, fmt.Errorf(ErrorWrongCredentials)
	}

	err = checkPassword(reqDto.Password, user.Password)
	if err != nil {
		return nil, err
	}

	token, err := getToken(user)
	if err != nil {
		return nil, err
	}

	resDto := responses.Login{
		Token: token,
	}

	return &resDto, nil
}

func checkPassword(plainPassword string, hashFromDb string) error {
	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash(plainPassword, hashFromDb)
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())
		return err
	}

	if match {
		return nil
	}

	return fmt.Errorf(ErrorWrongCredentials)
}
