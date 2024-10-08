package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	"sso/internal/storage"
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
	IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error)
}

type AppProvider interface {
	GetAppById(ctx context.Context, appID uuid.UUID) (models.App, error)
}

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

func NewAuth(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Register new user in the system and returns user ID
// If user with given username already exists, returns error
func (a *Auth) SignUp(ctx context.Context, username string, email string, password string) (uid uuid.UUID, err error) {
	const op = "Auth.SignUp"

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
func (a *Auth) SignIn(ctx context.Context, username string, password string) (string, error) {
	const op = "Auth.SignIn"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("attempting to sign in user")

	user, err := a.userProvider.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s, %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s, %w", op, err)
	}

	return token, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId string) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("userId", userId),
	)

	uid, err := uuid.Parse(userId)
	if err != nil {
		return false, fmt.Errorf("%s, %w", op, err)
	}

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, uid)
	if err != nil {
		return false, fmt.Errorf("%s, %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

func (a *Auth) UserIdentity(ctx context.Context, accessToken string) (bool, uuid.UUID, error) {
	const op = "Auth.UserIdentity"

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
