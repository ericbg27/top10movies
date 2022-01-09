package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	movies_controller "github.com/ericbg27/top10movies-api/src/controllers/movies"
	users_controller "github.com/ericbg27/top10movies-api/src/controllers/users"
	postgresdb "github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
	"github.com/ericbg27/top10movies-api/src/server"
	movies_service "github.com/ericbg27/top10movies-api/src/services/movies"
	users_service "github.com/ericbg27/top10movies-api/src/services/users"
	"github.com/ericbg27/top10movies-api/src/utils/authorization"
	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := startApplication(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func startApplication() error {
	cfg := config.GetConfig()

	logger.SetupLogger(cfg)

	router := gin.Default()

	db := postgresdb.NewPostgresDBClient(cfg.Database)

	db.SetupDbConnection()
	defer db.CloseDbConnection(context.Background())

	redisClient := &redisdb.RedisClient{}
	redisClient.SetupRedisConnection(cfg.Redis.CacheTtl)

	authorizationManager := authorization.NewAuthorizationManager(redisClient)

	usersService := users_service.NewUsersService(db, *redisClient)
	moviesService := movies_service.NewMoviesService(redisClient, cfg)

	usersController := users_controller.NewUsersController(usersService, moviesService, authorizationManager)
	moviesController := movies_controller.NewMoviesController(moviesService)

	s := &server.Server{
		Router:           router,
		UsersController:  usersController,
		MoviesController: moviesController,
	}

	var sb strings.Builder

	sb.WriteString(strings.TrimSpace(cfg.Server.Host))
	sb.WriteString(":")
	sb.WriteString(strings.TrimSpace(cfg.Server.Port))

	return s.StartServer(sb.String())
}
