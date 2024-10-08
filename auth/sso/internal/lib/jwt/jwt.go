package jwt

import (
	"sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	signingKey = "hwekjf#hadsujfDPDSFJO21adho@JDSOV*@79Q"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
