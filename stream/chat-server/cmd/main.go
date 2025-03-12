package main

import (
	"chat-server/internal/app"
	"chat-server/internal/config"
	"chat-server/internal/repository/postgres"
	"chat-server/internal/repository/redis"
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad(os.Getenv("CONFIG_PATH"))

	logger := setupLogger(cfg.Env)
	defer logger.Sync()

	sugar := logger.Sugar()

	ctx := context.Background()

	pdb := postgres.NewPostgresDB(ctx, cfg.DB)
	rdb := redis.NewRedisDB(ctx, cfg.RDB)

	application := app.NewApp(sugar, cfg.HTTP.Port, pdb, rdb)

	go func() {
		application.Server.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop(ctx)
	pdb.Close()
	sugar.Infow("gracefully stopped")
}

func setupLogger(env string) *zap.Logger {
	var logger *zap.Logger
	var err error

	switch env {
	case envLocal:
		logger = zap.NewExample()
	case envDev:
		logger, err = zap.NewDevelopment()
	case envProd:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to setup logger")
	}

	return logger
}
