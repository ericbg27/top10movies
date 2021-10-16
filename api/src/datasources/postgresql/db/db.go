package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	deleteCacheQuery = "DELETE FROM movies AS m1 WHERE EXTRACT(EPOCH FROM (NOW() - m1.created_at)) >= $1;"
)

var (
	Client *pgxpool.Conn

	host     = config.GetConfig().Database.Host
	port     = config.GetConfig().Database.Port
	user     = config.GetConfig().Database.User
	password = config.GetConfig().Database.Password
	dbname   = config.GetConfig().Database.DbName
	loglevel = config.GetConfig().Database.LogLevel
	cachettl = config.GetConfig().Database.CacheTtl
)

func SetupDbConnection() {
	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname))
	if err != nil {
		logger.Error("Error when parsing database connection string", err)
		panic(err)
	}

	config.ConnConfig.Logger = logger.GetLogger()

	level, err := pgx.LogLevelFromString(strings.ToLower(loglevel))
	if err != nil {
		level = pgx.LogLevelInfo
	}
	config.ConnConfig.LogLevel = level

	clientPool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to database: %v\n", err.Error()), err)
		panic(err)
	}

	Client, err = clientPool.Acquire(context.Background())
	if err != nil {
		logger.Error("Unable to acquire connection to database", err)
		panic(err)
	}
}

func ClearMoviesCache() {
	clearCacheTicker := time.NewTicker(time.Duration(cachettl) * time.Minute)

	for {
		<-clearCacheTicker.C

		logger.Info("Clearing movies cache")

		cachettlSeconds := cachettl * 60

		// TODO: Improve cache clearing by saving the timestamp when records are saved and verifying if the elapsed time is bigger than cachettl or not
		Client.Exec(context.Background(), deleteCacheQuery, cachettlSeconds)
	}
}
