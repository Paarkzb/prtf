package redis

import (
	"context"
	"fmt"
	"videostream/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisDB(ctx context.Context, cfg config.RDBConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return rdb
}
