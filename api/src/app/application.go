package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()

	db.SetupDbConnection()
	defer db.Client.Conn().Close(context.Background())

	redisdb.SetupRedisConnection()

	cfg := config.GetConfig()

	var sb strings.Builder

	sb.WriteString(strings.TrimSpace(cfg.Server.Host))
	sb.WriteString(":")
	sb.WriteString(strings.TrimSpace(cfg.Server.Port))

	logger.Info(fmt.Sprintf("Starting the application at %s", sb.String()))
	router.Run(sb.String())
}
