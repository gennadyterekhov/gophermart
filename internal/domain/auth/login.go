package auth

import (
	"context"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/gennadyterekhov/gophermart/internal/domain/auth/token"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/logger"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
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

	tokenString, err := token.CreateToken(user)
	if err != nil {
		return nil, err
	}

	resDto := responses.Login{
		Token: tokenString,
	}

	return &resDto, nil
}

func checkPassword(plainPassword string, hashFromDB string) error {
	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash(plainPassword, hashFromDB)
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())
		return err
	}

	if match {
		return nil
	}

	return fmt.Errorf(ErrorWrongCredentials)
}
