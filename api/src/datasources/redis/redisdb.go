package redisdb

import (
	"fmt"
	"os"

	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/go-redis/redis"
)

const (
	localRedis = "localhost:6379"
	RedisNil   = redis.Nil
)

type RedisClient struct {
	Client   *redis.Client
	CacheTTL int64
}

func (r *RedisClient) SetupRedisConnection(cacheTtl int64) {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = localRedis
	}

	r.Client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	_, err := r.Client.Ping().Result()
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to Redis at %s: %s\n", dsn, err.Error()), err)

		panic(err)
	}

	r.CacheTTL = cacheTtl

	logger.Info(fmt.Sprintf("Connected to Redis at %s", dsn))
}
