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
	users_service "github.com/ericbg27/top10movies-api/src/services/users"
	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := startApplication(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func startApplication() error {
	router := gin.Default()

	db := &postgresdb.PostgresDBClient{
		Client: nil,
	}

	db.SetupDbConnection()
	defer db.CloseDbConnection(context.Background())

	users_service.UsersService.SetupDBClient(db)

	usersController := users_controller.NewUsersController()
	moviesController := movies_controller.NewMoviesController()

	redisdb.SetupRedisConnection()

	s := &server.Server{
		Db:               db,
		Router:           router,
		UsersController:  usersController,
		MoviesController: moviesController,
	}

	cfg := config.GetConfig()

	var sb strings.Builder

	sb.WriteString(strings.TrimSpace(cfg.Server.Host))
	sb.WriteString(":")
	sb.WriteString(strings.TrimSpace(cfg.Server.Port))

	return s.StartApplication(sb.String())
}
