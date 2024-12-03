package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	"sso/internal/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSaver interface {
	SaveUser(ctx context.Context, username string, email string, passHash []byte) (uid uuid.UUID, err error)
}

type UserProvider interface {
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	GetUserByUserID(ctx context.Context, userID uuid.UUID) (models.User, error)
}

type TokenProvider interface {
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, refreshTokenTTL time.Duration) error
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, refreshTokenTTL time.Duration) error
	GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
}

type AppProvider interface {
	GetAppById(ctx context.Context, appID uuid.UUID) (models.App, error)
}

type AuthService struct {
	log             *slog.Logger
	userSaver       UserSaver
	userProvider    UserProvider
	appProvider     AppProvider
	tokenProvider   TokenProvider
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenProvider TokenProvider, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:             log,
		userSaver:       userSaver,
		userProvider:    userProvider,
		appProvider:     appProvider,
		tokenProvider:   tokenProvider,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// Register new user in the system and returns user ID
// If user with given username already exists, returns error
func (a *AuthService) SignUp(ctx context.Context, username string, email string, password string) (uid uuid.UUID, err error) {
	const op = "AuthService.SignUp"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("Register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, username, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))

		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}

// SignIn checks if user with given credentials exists in the system and returns access token
// If user exists and password is incorrect, returns error
// If user doesn't exist, returns error
func (a *AuthService) SignIn(ctx context.Context, username string, password string) (models.Tokens, error) {
	const op = "AuthService.SignIn"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("attempting to sign in user")

	var tokens models.Tokens

	user, err := a.userProvider.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return tokens, fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return tokens, fmt.Errorf("%s, %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return tokens, fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
	}

	tokens, err = generateTokens(user, a.accessTokenTTL)
	if err != nil {
		a.log.Info("failed to generate tokens", sl.Err(err))

		return tokens, fmt.Errorf("%s, %w", op, err)
	}

	err = a.tokenProvider.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken, a.refreshTokenTTL)
	if err != nil {
		a.log.Info("failed to save refresh token", sl.Err(err))

		return tokens, fmt.Errorf("%s, %w", op, err)
	}

	return tokens, nil
}

func generateTokens(user models.User, accessTokenTTL time.Duration) (models.Tokens, error) {
	var tokens models.Tokens

	accessToken, err := jwt.NewAccessToken(user, accessTokenTTL)
	if err != nil {

		return tokens, err
	}

	refreshToken, err := jwt.NewRefreshToken()
	if err != nil {

		return tokens, err
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "AuthService.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("userID", userID.String()),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s, %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

func (a *AuthService) UserIdentity(ctx context.Context, accessToken string) (bool, uuid.UUID, error) {
	const op = "AuthService.userIdentity"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("authenticate user")
	claims, err := jwt.ParseToken(accessToken)
	if err != nil {
		return false, uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	uid, err := uuid.Parse(claims["uid"].(string))
	if err != nil {
		return false, uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}
	log.Info("user authenticated", slog.String("user_id", uid.String()))

	return true, uid, nil
}

func (a *AuthService) Refresh(ctx context.Context, userID uuid.UUID, refreshToken string) (models.Tokens, error) {
	const op = "AuthService.Refresh"
	var tokens models.Tokens

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("refresh tokens")
	dbRefreshToken, err := a.tokenProvider.GetRefreshToken(ctx, userID)
	if err != nil {
		a.log.Error("user session not found", sl.Err(err))

		return tokens, fmt.Errorf("%s, %w", op, err)
	}

	if dbRefreshToken != refreshToken {
		a.log.Error("invalid refresh token")

		return tokens, fmt.Errorf("%s, %w", op, errors.New("invalid refresh token"))
	}

	user, err := a.userProvider.GetUserByUserID(ctx, userID)
	if err != nil {
		a.log.Error("user not found")

		return tokens, fmt.Errorf("%s, %w", op, err)
	}

	tokens, err = generateTokens(user, a.accessTokenTTL)
	if err != nil {
		a.log.Info("failed to generate tokens", sl.Err(err))

		return tokens, fmt.Errorf("%s, %w", op, err)
	}

	err = a.tokenProvider.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken, a.refreshTokenTTL)
	if err != nil {
		a.log.Info("failed to save refresh token", sl.Err(err))

		return tokens, fmt.Errorf("%s, %w", op, err)
	}

	return tokens, nil

}
