package tests

import (
	ssov1 "sso/protos/gen/go/sso"
	"sso/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appID     = "36c604ca-5f22-447c-a2a7-f220d2c1193b"
	appSecret = "test-secret"
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	const passDefaultLen = 10

	username := gofakeit.Username()
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	respReg, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
		Username: username,
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)

	// TODO: Проверить результат
}
