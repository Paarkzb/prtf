package repository

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserExists          = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserSessionNotFound = errors.New("user session not found")
	ErrAppNotFound         = errors.New("app not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) SaveUser(ctx context.Context, username string, email string, passHash []byte) (uid uuid.UUID, err error) {
	const op = "Repository.postgres.SaveUser"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	query := `INSERT INTO public.users(username, email, pass_hash) VALUES($1, $2, $3) RETURNING id`
	var id uuid.UUID

	err = tx.QueryRow(ctx, query, username, email, passHash).Scan(&id)
	if err != nil {
		_ = tx.Rollback(ctx)

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, tx.Commit(ctx)
}

func (s *Repository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	const op = "Repository.postgres.GetUserByUsername"

	var user models.User

	query := `SELECT id, username, email, pass_hash FROM public.users WHERE username=$1`

	err := s.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Repository) GetAppById(ctx context.Context, appID uuid.UUID) (models.App, error) {
	const op = "Repository.postgres.GetUserByUsername"

	var app models.App

	query := `SELECT id, name, secret FROM public.apps WHERE id=$1`

	err := s.db.QueryRow(ctx, query, appID).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app, fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}
		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Repository) IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error) {
	const op = "Repository.postgres.IsAdmin"

	query := `SELECT id FROM public.users_admins WHERE rf_users_id=$1`

	var uid uuid.UUID
	err := s.db.QueryRow(ctx, query, userId).Scan(&uid)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	isAdmin := true
	if uid == uuid.Nil {
		isAdmin = false
	}

	return isAdmin, nil
}

func (s *Repository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, refreshTokenTTL time.Duration) error {
	const op = "Repository.postgres.SaveRefreshToken"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `INSERT INTO public.users_sessions(rf_users_id, refresh_token, expires_at) VALUES ($1, $2, $3) 
	ON CONFLICT (rf_users_id) DO UPDATE SET refresh_token=$2, expires_at=$3 WHERE users_sessions.rf_users_id=$1`

	expTime := time.Now().Add(refreshTokenTTL)
	_, err = tx.Exec(ctx, query, userID, refreshToken, expTime)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}

func (s *Repository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, refreshTokenTTL time.Duration) error {
	const op = "Repository.postgres.SaveRefreshToken"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `UPDATE public.users_sessions SET refresh_token=$1, expires_at=$2 WHERE rf_users_id=$3`

	expTime := time.Now().Add(refreshTokenTTL)
	_, err = tx.Exec(ctx, query, refreshToken, expTime, userID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}

func (s *Repository) GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	const op = "Repository.postgres.GetRefreshToken"

	query := `SELECT refresh_token FROM public.users_sessions WHERE rf_users_id=$1`

	var refreshToken string
	err := s.db.QueryRow(ctx, query, userID).Scan(&refreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, ErrUserSessionNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, nil
}

func (s *Repository) GetUserByUserID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	const op = "Repository.postgres.GetUserByUserID"

	var user models.User

	query := `SELECT id, username, email FROM public.users WHERE id=$1`

	err := s.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
