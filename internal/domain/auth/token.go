package auth

import (
	"fmt"
	"os"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func getToken(user *models.User) (string, error) {
	var (
		token         *jwt.Token
		tokenAsString string
		err           error
	)

	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "gophermart",
			"sub": user.Login,
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

func validateToken(token string, login string) error {
	claims, err := getClaimsFromToken(token)
	if err != nil {
		return err
	}
	sub, err := claims.GetSubject()
	if err != nil {
		return err
	}

	if sub == login {
		return nil
	}

	return fmt.Errorf("token did not authenticate selected user")
}

func getClaimsFromToken(token string) (*jwt.MapClaims, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	tokenObject, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return getJwtSigningKey(), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "error when parsing token")
	}
	claims, ok := tokenObject.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error when getting claims from token")
	}
	return &claims, nil
}

func GetLoginFromToken(token string) (string, error) {
	claims, err := getClaimsFromToken(token)
	if err != nil {
		return "", err
	}

	login, err := claims.GetSubject()
	if err != nil {
		return "", err
	}

	return login, nil
}
