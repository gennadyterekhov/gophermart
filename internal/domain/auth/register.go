package auth

import (
	"context"
	"os"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"

	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"

	"github.com/alexedwards/argon2id"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/golang-jwt/jwt/v5"
)

func Register(ctx context.Context, reqDto *requests.Register) (*responses.Register, error) {
	encryptedPassword, err := encrypt(reqDto.Password)
	if err != nil {
		return nil, err
	}

	user, err := repositories.AddUser(ctx, reqDto.Login, encryptedPassword)
	if err != nil {
		return nil, err
	}

	token, err := getToken(user)
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

func getToken(user *models.User) (string, error) {
	var (
		token         *jwt.Token
		tokenAsString string
		err           error
	)

	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "gophermart",
			"sub": user.ID,
		},
	)

	tokenAsString, err = token.SignedString(getJwtSigningKey())
	if err != nil {
		return "", err
	}

	return tokenAsString, nil
}

func getJwtSigningKey() []byte {
	fromEnv, ok := os.LookupEnv("JWT_SIGNING_KEY")
	if ok {
		return []byte(fromEnv)
	}

	return []byte("")
}
