package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

const (
	userTable     = "public.user"
	quizTable     = "public.quiz"
	questionTable = "public.question"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg DBConfig) (*pgxpool.Pool, error) {
	// "postgres://username:password@localhost:5432/database_name"
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
