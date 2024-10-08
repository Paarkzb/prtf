package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const (
	signingKey = "hwekjf#hadsujfDPDSFJO21adho@JDSOV*@79Q"
)

func ParseToken(accessToken string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims, nil
}
