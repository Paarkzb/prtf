package repository

import (
	"context"
	"fmt"
	"time"
	"videostream/internal/domain/models"

	"github.com/redis/go-redis/v9"
)

type RepositoryRedis struct {
	rdb *redis.Client
}

func NewRepositoryRedis(rdb *redis.Client) *RepositoryRedis {
	return &RepositoryRedis{
		rdb: rdb,
	}
}

func (r *RepositoryRedis) SetChannels(ctx context.Context, channels []models.Channel) error {
	const op = "Repository.redis.SetChannels"

	err := r.rdb.Set(ctx, "channels", channels, 10*time.Second).Err()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RepositoryRedis) GetChannels(ctx context.Context) ([]models.Channel, error) {
	const op = "Repository.redis.GetChannels"

	var channels []models.Channel

	err := r.rdb.Get(ctx, "channels").Scan(&channels)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return channels, nil
}
