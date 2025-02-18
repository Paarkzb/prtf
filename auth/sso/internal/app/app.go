package app

import (
	"context"
	"sso/internal/app/server"
	"sso/internal/handler"
	"sso/internal/repository"
	"sso/internal/services/authservice"
	"time"

	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	HTTPServer *server.App
}

func NewApp(log *slog.Logger, port int, db *pgxpool.Pool, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *App {

	authRepo := repository.NewRepository(db)
	authService := authservice.NewAuthService(log, authRepo, authRepo, authRepo, authRepo, accessTokenTTL, refreshTokenTTL)
	handlers := handler.NewHandler(log, authService)

	httpServer := server.NewApp(log, handlers.InitRoutes(), port)

	return &App{
		HTTPServer: httpServer,
	}
}

func (a *App) Stop(ctx context.Context) error {

	return a.HTTPServer.Stop(ctx)
}
