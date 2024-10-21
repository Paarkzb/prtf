package app

import (
	"context"
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	httpgateway "sso/internal/app/http-gateway"
	authservice "sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
	HTTPServer *httpgateway.App
}

func NewApp(ctx context.Context, log *slog.Logger, grpcPort int, httpPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.NewStorage(ctx, storagePath)
	if err != nil {
		panic(err)
	}

	authService := authservice.NewAuth(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewApp(log, authService, grpcPort)

	httpApp := httpgateway.NewApp(log, httpPort)

	return &App{
		GRPCServer: grpcApp,
		HTTPServer: httpApp,
	}
}
