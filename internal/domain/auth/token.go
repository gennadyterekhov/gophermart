package auth

import (
	"fmt"
	"os"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
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

func validateToken(token string, id int64) error {
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
		return fmt.Errorf("error when parsing token")
	}
	claims, ok := tokenObject.Claims.(jwt.MapClaims)

	if ok {
		if claims["sub"] == float64(id) {
			return nil
		}
	}
	return fmt.Errorf("token did not authenticate selected user")
}
