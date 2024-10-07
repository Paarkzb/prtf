package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	signingKey = "hwekjf#hadsujfDPDSFJO21adho@JDSOV*@79Q"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId uuid.UUID `json:"user_id"`
}

func ParseToken(accessToken string) (tokenClaims, error) {
	var claims tokenClaims

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return claims, err
	}

	claims, ok := token.Claims.(tokenClaims)
	if !ok {
		return claims, errors.New("token claims are not of type *tokenClaims")
	}

	return claims, nil
}
