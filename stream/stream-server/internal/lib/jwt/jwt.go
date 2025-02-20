package jwt

import (
	"errors"
	"time"
	"videostream/internal/domain/models"

	"github.com/golang-jwt/jwt"
)

func ParseToken(accessToken string) (jwt.MapClaims, error) {

	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot cast claims")
	}

	if claims["uid"] == "" {
		return nil, errors.New("invalid token payload")
	}

	return claims, nil
}

const (
	signingKey = "laldsnvopPO@U!@POJfopjx?><M@!KM/lvkdsj"
)

func NewStreamToken(channel models.Channel, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = channel.RfUserID
	claims["channel_token"] = channel.ChannelToken
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["iat"] = time.Now().Unix()

	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ParseStreamToken(streamToken string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(streamToken, func(t *jwt.Token) (interface{}, error) {
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
