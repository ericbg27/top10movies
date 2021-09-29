package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/jackc/pgx"
)

const (
	deleteCacheQuery = "DELETE FROM movies;"
)

var (
	Client *pgx.ConnPool

	host     = config.GetConfig().Database.Host
	port     = config.GetConfig().Database.Port
	user     = config.GetConfig().Database.User
	password = config.GetConfig().Database.Password
	dbname   = config.GetConfig().Database.DbName
	loglevel = config.GetConfig().Database.LogLevel
	cachettl = config.GetConfig().Database.CacheTtl
)

func SetupDbConnection() {
	var connconfig pgx.ConnPoolConfig
	connconfig.Host = host
	connconfig.Port = port
	connconfig.Password = password
	connconfig.Database = dbname
	connconfig.User = user

	connconfig.Logger = logger.GetLogger()

	level, err := pgx.LogLevelFromString(strings.ToLower(loglevel))
	if err != nil {
		level = pgx.LogLevelInfo
	}
	connconfig.LogLevel = level

	Client, err = pgx.NewConnPool(connconfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to database: %v\n", err.Error()), err)
		panic(err)
	}
}

func ClearMoviesCache() {
	clearCacheTicker := time.NewTicker(time.Duration(cachettl) * time.Second)

	for {
		<-clearCacheTicker.C

		logger.Info("Clearing movies cache")

		// TODO: Improve cache clearing by saving the timestamp when records are saved and verifying if the elapsed time is bigger than cachettl or not
		Client.Exec(deleteCacheQuery)
	}
}
