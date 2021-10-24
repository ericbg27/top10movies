package redisdb

import (
	"fmt"
	"os"

	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/go-redis/redis"
)

const (
	localRedis = "localhost:6379"
)

var (
	Client *redis.Client
)

func SetupRedisConnection() {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = localRedis
	}

	Client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	_, err := Client.Ping().Result()
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to Redis at %s: %s\n", dsn, err.Error()), err)

		panic(err)
	}

	logger.Info(fmt.Sprintf("Connected to Redis at %s", dsn))
}
