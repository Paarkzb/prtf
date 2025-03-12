package repository

import "github.com/redis/go-redis/v9"

type RepositoryRedis struct {
	rdb *redis.Client
}

func NewRepositoryRedis(rdb *redis.Client) *RepositoryRedis {
	return &RepositoryRedis{
		rdb: rdb,
	}
}
