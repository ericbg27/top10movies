package server

import (
	"fmt"

	"github.com/gin-gonic/gin"

	movies_controller "github.com/ericbg27/top10movies-api/src/controllers/movies"
	users_controller "github.com/ericbg27/top10movies-api/src/controllers/users"
	"github.com/ericbg27/top10movies-api/src/datasources/database"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
)

type Server struct {
	Db               database.DatabaseClient
	Router           *gin.Engine
	UsersController  users_controller.UsersControllerInterface
	MoviesController movies_controller.MoviesControllerInterface
}

func (s *Server) StartApplication(address string) error {
	s.routes()

	logger.Info(fmt.Sprintf("Starting the application at %s", address))

	if err := s.Router.Run(address); err != nil {
		return err
	}

	return nil
}
