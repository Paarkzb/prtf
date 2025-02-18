package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	log        *zap.SugaredLogger
	httpServer *http.Server
}

func NewServer(log *zap.SugaredLogger, handler http.Handler, port int) *Server {
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return &Server{
		log:        log,
		httpServer: httpServer,
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Run() error {
	s.log.Infow("http server started ", "Addr: ", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}
