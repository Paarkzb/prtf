package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"videostream/internal/app"
	"videostream/internal/config"
	"videostream/internal/repository/postgres"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "prtf-stream-redis:6379",
	Password: "redis",
	DB:       0,
})

func recordMetrics() {
	go func() {
		for {
			// promActiveStreams.Add(float64(len(activeStreams)))
			promActiveStreams.Add(23)

			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	promStreamRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "stream_requests_total",
			Help: "Total number of stream requests",
		},
	)
	promActiveStreams = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_streams",
			Help: "Current active streams",
		},
	)
)

func init() {
	prometheus.MustRegister(promStreamRequests, promActiveStreams)
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	recordMetrics()

	cfg := config.MustLoad(os.Getenv("CONFIG_PATH"))

	logger := setupLogger(cfg.Env)
	defer logger.Sync()

	sugar := logger.Sugar()

	ctx := context.Background()

	db, err := postgres.NewPostgresDB(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}

	// redis
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}

	application := app.NewApp(sugar, cfg.HTTP.Port, db)

	go func() {
		application.Server.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop(ctx)
	db.Close()
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
