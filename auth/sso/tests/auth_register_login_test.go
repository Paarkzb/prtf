package tests

import (
	ssov1 "sso/protos/gen/go/sso"
	"sso/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appID     = "36c604ca-5f22-447c-a2a7-f220d2c1193b"
	appSecret = "hwekjfskladjvhiweuhfwieuh"
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	const passDefaultLen = 10

	username := gofakeit.Username()
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	respReg, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
		Email:    email,
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), claims["uid"].(uuid.UUID))
	assert.Equal(t, username, claims["username"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}
