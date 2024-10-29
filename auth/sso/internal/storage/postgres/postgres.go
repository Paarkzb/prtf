package postgres

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, username string, email string, passHash []byte) (uid uuid.UUID, err error) {
	const op = "storage.postgres.SaveUser"

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

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	const op = "storage.postgres.GetUserByUsername"

	var user models.User

	query := `SELECT id, username, email, pass_hash FROM public.users WHERE username=$1`

	err := s.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetAppById(ctx context.Context, appID uuid.UUID) (models.App, error) {
	const op = "storage.postgres.GetUserByUsername"

	var app models.App

	query := `SELECT id, name, secret FROM public.apps WHERE id=$1`

	err := s.db.QueryRow(ctx, query, appID).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	query := `SELECT id FROM public.users_admins WHERE rf_users_id=$1`

	var uid uuid.UUID
	err := s.db.QueryRow(ctx, query, userId).Scan(&uid)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	isAdmin := true
	if uid == uuid.Nil {
		isAdmin = false
	}

	return isAdmin, nil
}

func (s *Storage) SaveRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, refreshTokenTTL time.Duration) error {
	const op = "storage.postgres.SaveRefreshToken"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `SELECT id FROM public.users_sessions as us WHERE rf_users_id=$1`
	var id uuid.UUID

	err = tx.QueryRow(ctx, query, userID).Scan(&id)
	if err != nil {
		_ = tx.Rollback(ctx)

		return fmt.Errorf("%s: %w", op, err)
	}

	query = `UPDATE public.users_sessions SET refresh_token=$1, expires_at=$2 WHERE id=$3`

	expTime := time.Now().Add(refreshTokenTTL)
	_, err = tx.Exec(ctx, query, refreshToken, expTime, id)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}
