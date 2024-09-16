package app

import (
	"context"
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	authservice "sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(ctx context.Context, log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.NewStorage(ctx, storagePath)
	if err != nil {
		panic(err)
	}

	authService := authservice.NewAuth(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewApp(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
