package jwt

import (
	"errors"

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
