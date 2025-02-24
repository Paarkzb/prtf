package postgres

import (
	"context"
	"fmt"
	"videostream/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg config.DBConfig) *pgxpool.Pool {
	// "postgres://username:password@localhost:5432/database_name"
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))
	if err != nil {
		panic(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return pool
}
