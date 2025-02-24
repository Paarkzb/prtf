package app

import (
	"context"
	"videostream/internal/app/server"
	"videostream/internal/handler"
	"videostream/internal/repository"
	"videostream/internal/services/streamservice"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type App struct {
	Server *server.Server
}

func NewApp(log *zap.SugaredLogger, port int, pdb *pgxpool.Pool, rdb *redis.Client) *App {

	streamRepo := repository.NewRepositoryPostgres(pdb)
	redisRepo := repository.NewRepositoryRedis(rdb)

	streamService := streamservice.NewStreamService(log, streamRepo, streamRepo, redisRepo)

	streamHandler := handler.NewHandler(log, streamService)

	httpServer := server.NewServer(log, streamHandler.InitRoutes(), port)

	return &App{
		Server: httpServer,
	}
}

func (a *App) Stop(ctx context.Context) error {

	return a.Server.Stop(ctx)
}
