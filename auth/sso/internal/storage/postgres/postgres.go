package postgres

import (
	"context"
	"fmt"
	"sso/internal/domain/models"

	"github.com/google/uuid"
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

	return &Storage{db: db}, err
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

}
