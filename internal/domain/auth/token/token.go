package token

import (
	"fmt"
	"os"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/domain/auth/jwtclaims"
	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const issuerGophermart = "gophermart"

func CreateToken(user *models.User) (string, error) {
	var (
		token         *jwt.Token
		tokenAsString string
		err           error
	)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		&jwtclaims.Claims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().AddDate(1, 0, 0)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			Issuer:    issuerGophermart,
			Subject:   user.Login,
			Audience:  jwt.ClaimStrings{},
			UserID:    user.ID,
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

func ValidateToken(token string, login string) error {
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

func getClaimsFromToken(token string) (*jwtclaims.Claims, error) {
	claims := &jwtclaims.Claims{}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	_, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return getJwtSigningKey(), nil
		},
	)
	if err != nil {
		logger.CustomLogger.Errorln("could not parse token ", token)
		return nil, errors.Wrap(err, "error when parsing token")
	}

	return claims, nil
}

func GetIDAndLoginFromToken(token string) (int64, string, error) {
	claims, err := getClaimsFromToken(token)
	if err != nil {
		return 0, "", err
	}

	login, err := claims.GetSubject()
	if err != nil {
		return 0, "", err
	}
	id, err := claims.GetUserID()
	if err != nil {
		return 0, "", err
	}

	return id, login, nil
}
