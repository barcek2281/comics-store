package jwt

import (
	"time"

	"github.com/barcek2281/comics-store/auth/internal/model"
	"github.com/golang-jwt/jwt"
)

func NewToken(secret string, user model.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
