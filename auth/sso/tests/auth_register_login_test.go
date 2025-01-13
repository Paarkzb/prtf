package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sso/internal/domain/models"
	"sso/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appID          = "36c604ca-5f22-447c-a2a7-f220d2c1193b"
	appSecret      = "hwekjf#hadsujfDPDSFJO21adho@JDSOV*@79Q"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	st := suite.NewSuite(t)

	username := gofakeit.Username()
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	userInput := models.UserInput{
		Email:    email,
		Username: username,
		Password: pass,
	}
	data, err := json.Marshal(userInput)
	require.NoError(t, err)
	r := bytes.NewReader(data)
	respReg, err := st.AuthClient.Post(fmt.Sprintf("%s:%d/v1/sign-up", "http://localhost", st.Cfg.HTTP.Port), "application/json", r)
	require.NoError(t, err)

	regBody, err := io.ReadAll(respReg.Body)
	require.NoError(t, err)
	var regData map[string]string
	err = json.Unmarshal(regBody, &regData)
	require.NoError(t, err)
	assert.NotEmpty(t, regData["userID"])

	userInput = models.UserInput{
		Username: username,
		Password: pass,
	}
	data, err = json.Marshal(userInput)
	require.NoError(t, err)
	r = bytes.NewReader(data)
	respLogin, err := st.AuthClient.Post(fmt.Sprintf("%s:%d/v1/sign-in", "http://localhost", st.Cfg.HTTP.Port), "application/json", r)
	require.NoError(t, err)

	loginBody, err := io.ReadAll(respLogin.Body)
	require.NoError(t, err)
	var loginData map[string]map[string]string
	err = json.Unmarshal(loginBody, &loginData)
	require.NoError(t, err)

	token := loginData["tokens"]["access_token"]
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, regData["userID"], claims["uid"].(string))
	assert.Equal(t, username, claims["username"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.AccessTokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	st := suite.NewSuite(t)

	username := gofakeit.Username()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	userInput := models.UserInput{
		Email:    email,
		Username: username,
		Password: password,
	}
	data, err := json.Marshal(userInput)
	require.NoError(t, err)
	r := bytes.NewReader(data)
	respReg, err := st.AuthClient.Post(fmt.Sprintf("%s:%d/v1/sign-up", "http://localhost", st.Cfg.HTTP.Port), "application/json", r)
	require.NoError(t, err)

	regBody, err := io.ReadAll(respReg.Body)
	require.NoError(t, err)
	var regData map[string]string
	err = json.Unmarshal(regBody, &regData)
	require.NoError(t, err)
	assert.NotEmpty(t, regData["userID"])

	data, err = json.Marshal(userInput)
	require.NoError(t, err)
	r = bytes.NewReader(data)
	respReg, err = st.AuthClient.Post(fmt.Sprintf("%s:%d/v1/sign-up", "http://localhost", st.Cfg.HTTP.Port), "application/json", r)
	require.NoError(t, err)

	regBody, err = io.ReadAll(respReg.Body)
	require.NoError(t, err)
	var regData2 map[string]string
	err = json.Unmarshal(regBody, &regData2)
	require.NoError(t, err)

	assert.Empty(t, regData2["userID"])
	require.Contains(t, regData2["message"], "failed to register user")
}

func TestRegister_FailCases(t *testing.T) {
	st := suite.NewSuite(t)

	tests := []struct {
		name        string
		username    string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "SignUp with empty password",
			username:    gofakeit.Username(),
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "SignUp with empty username",
			username:    "",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			expectedErr: "username is required",
		},
		{
			name:        "SignUp with empty both",
			username:    "",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "username is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userInput := models.UserInput{
				Email:    tt.email,
				Username: tt.username,
				Password: tt.password,
			}
			data, err := json.Marshal(userInput)
			require.NoError(t, err)
			r := bytes.NewReader(data)
			respReg, err := st.AuthClient.Post(fmt.Sprintf("%s:%d/v1/sign-up", "http://localhost", st.Cfg.HTTP.Port), "application/json", r)
			require.NoError(t, err)

			regBody, err := io.ReadAll(respReg.Body)
			require.NoError(t, err)
			var regData map[string]string
			err = json.Unmarshal(regBody, &regData)
			require.NoError(t, err)

			require.Contains(t, regData["message"], tt.expectedErr)
		})
	}
}

// func TestLogin_FailCases(t *testing.T) {
// 	ctx, st := suite.NewSuite(t)

// 	tests := []struct {
// 		name        string
// 		username    string
// 		email       string
// 		password    string
// 		expectedErr string
// 	}{
// 		{
// 			name:        "SignIn with empty password",
// 			username:    gofakeit.Username(),
// 			password:    "",
// 			expectedErr: "password is required",
// 		},
// 		{
// 			name:        "SignIn with empty username",
// 			username:    "",
// 			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
// 			expectedErr: "username is required",
// 		},
// 		{
// 			name:        "SignIn with empty both",
// 			username:    "",
// 			password:    "",
// 			expectedErr: "username is required",
// 		},
// 		{
// 			name:        "SignIn with not-matching username or password",
// 			username:    gofakeit.Username(),
// 			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
// 			expectedErr: "invalid username or password",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
// 				Username: gofakeit.Username(),
// 				Email:    gofakeit.Email(),
// 				Password: gofakeit.Password(true, true, true, true, false, passDefaultLen),
// 			})
// 			require.NoError(t, err)

// 			_, err = st.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
// 				Username: tt.username,
// 				Password: tt.password,
// 			})
// 			require.Error(t, err)
// 			require.Contains(t, err.Error(), tt.expectedErr)
// 		})
// 	}
// }
