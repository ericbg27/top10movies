package users_db

import (
	"fmt"
	"strings"

	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/jackc/pgx"
)

var (
	Client *pgx.ConnPool

	host     = config.GetConfig().Database.Host
	port     = config.GetConfig().Database.Port
	user     = config.GetConfig().Database.User
	password = config.GetConfig().Database.Password
	dbname   = config.GetConfig().Database.DbName
	loglevel = config.GetConfig().Database.LogLevel
)

/*func init() {
	setupDbConnection()
}*/

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
