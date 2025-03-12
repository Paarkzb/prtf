package app

import (
	"chat-server/internal/app/server"
	"chat-server/internal/handler"
	"chat-server/internal/repository"
	"chat-server/internal/services/websocketservice"
	"context"

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

	// chatService := chatservice.NewChatService(log, streamRepo, redisRepo)

	websocketService := websocketservice.NewWebsocketService(log, streamRepo, redisRepo)

	chatHandler := handler.NewHandler(log, websocketService)

	httpServer := server.NewServer(log, chatHandler.InitRoutes(), port)

	return &App{
		Server: httpServer,
	}
}

func (a *App) Stop(ctx context.Context) error {

	return a.Server.Stop(ctx)
}
