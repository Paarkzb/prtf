package jwt

import (
	"errors"
	"math/rand"
	"sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey = "hwekjf#hadsujfDPDSFJO21adho@JDSOV*@79Q"
)

func NewAccessToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["iat"] = time.Now().Unix()

	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ParseToken(accessToken string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot cast claims")
	}

	return claims, nil
}

func NewRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	r := rand.New(rand.NewSource(time.Now().Unix()))

	if _, err := r.Read(bytes); err != nil {
		return "", err
	}

	refreshToken, err := bcrypt.GenerateFromPassword(bytes, 10)
	if err != nil {
		return "", err
	}

	return string(refreshToken), err
}
