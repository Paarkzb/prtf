package httpgateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	ssov1 "sso/protos/gen/go/sso"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log        *slog.Logger
	hTTPServer *http.Server
	port       int
}

func NewApp(log *slog.Logger, port int) *App {
	const op = "http-gateway.app.NewApp"

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	mux := runtime.NewServeMux()

	err := ssov1.RegisterAuthHandlerFromEndpoint(context.Background(), mux, "localhost:8085", opts)
	if err != nil {
		panic(err)
	}

	hTTPServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &App{
		log:        log,
		hTTPServer: hTTPServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "http-gateway.app.Run"

	a.log.Info("http server started", slog.String("addr", a.hTTPServer.Addr))

	return a.hTTPServer.ListenAndServe()
}
