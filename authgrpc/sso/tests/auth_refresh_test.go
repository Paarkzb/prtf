package tests

import (
	ssov1 "sso/protos/gen/go/sso"
	"sso/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefresh_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	username := gofakeit.Username()
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	respReg, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
		Email:    email,
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)

	respLogin, err := st.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)

	respRefresh, err := st.AuthClient.Refresh(ctx, &ssov1.RefreshRequest{
		UserID:       respReg.GetUserID(),
		RefreshToken: respLogin.GetRefreshToken(),
	})
	require.NoError(t, err)

	assert.Equal(t, respLogin.GetAccessToken(), respRefresh.GetAccessToken())
	assert.NotEqual(t, respLogin.GetRefreshToken(), respRefresh.GetRefreshToken())
}

func TestRefresh_FailCases(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	username := gofakeit.Username()
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	respReg, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
		Username: username,
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)

	respLogin, err := st.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)

	tests := []struct {
		name         string
		userID       string
		refreshToken string
		expectedErr  string
	}{
		{
			name:         "Refresh with empty userID",
			userID:       "",
			refreshToken: respLogin.GetRefreshToken(),
			expectedErr:  "userID is required",
		},
		{
			name:         "Refresh with empty refreshToken",
			userID:       respReg.GetUserID(),
			refreshToken: "",
			expectedErr:  "refresh token is required",
		},
		{
			name:         "Refresh with empty both",
			userID:       "",
			refreshToken: "",
			expectedErr:  "userID is required",
		},
		{
			name:         "Refresh with not-matching userID",
			userID:       gofakeit.UUID(),
			refreshToken: respLogin.GetRefreshToken(),
			expectedErr:  "failed to refresh token",
		},
		{
			name:         "Refresh with not-matching refreshToken",
			userID:       gofakeit.UUID(),
			refreshToken: "not-matching refreshToken",
			expectedErr:  "failed to refresh token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			respRefresh, err := st.AuthClient.Refresh(ctx, &ssov1.RefreshRequest{
				UserID:       tt.userID,
				RefreshToken: tt.refreshToken,
			})
			require.Error(t, err)

			assert.Empty(t, respRefresh.GetAccessToken(), respRefresh.GetRefreshToken())

			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
