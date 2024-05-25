package auth

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

const ErrorNotUniqueLogin = "ERROR: duplicate key value violates unique constraint \"users_login_key\" (SQLSTATE 23505)"

func Register(ctx context.Context, reqDto *requests.Register) (*responses.Register, error) {
	encryptedPassword, err := encrypt(reqDto.Password)
	if err != nil {
		return nil, err
	}

	user, err := repositories.AddUser(ctx, reqDto.Login, encryptedPassword)
	if err != nil {
		return nil, err
	}

	token, err := createToken(user)
	if err != nil {
		return nil, err
	}

	resDto := responses.Register{
		ID:    user.ID,
		Token: token,
	}

	return &resDto, nil
}

func encrypt(plainPassword string) (string, error) {
	// CreateHash returns an Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	hash, err := argon2id.CreateHash(plainPassword, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, err
}
