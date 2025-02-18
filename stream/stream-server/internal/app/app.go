package app

import (
	"context"
	"videostream/internal/app/server"
	"videostream/internal/handler"
	"videostream/internal/repository"
	"videostream/internal/services/streamservice"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	Server *server.Server
}

func NewApp(log *zap.SugaredLogger, port int, db *pgxpool.Pool) *App {

	streamRepo := repository.NewRepositoryPostgres(db)
	streamService := streamservice.NewStreamService(log, streamRepo)
	streamHandler := handler.NewHandler(log, streamService)

	httpServer := server.NewServer(log, streamHandler.InitRoutes(), port)
	return &App{
		Server: httpServer,
	}
}

func (a *App) Stop(ctx context.Context) error {

	return a.Server.Stop(ctx)
}
