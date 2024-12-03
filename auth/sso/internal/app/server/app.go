package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
}

func NewApp(log *slog.Logger, handler http.Handler, port int) *App {
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return &App{
		log:        log,
		httpServer: httpServer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.log.Info("http server started", slog.String("addr", a.httpServer.Addr))

	return a.httpServer.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error{
	return a.httpServer.Shutdown(ctx)
}
