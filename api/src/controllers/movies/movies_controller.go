package movies

import (
	"net/http"
	"strings"

	movies_service "github.com/ericbg27/top10movies-api/src/services/movies"
	"github.com/gin-gonic/gin"
)

type moviesController struct {
	moviesService movies_service.MoviesServiceInterface
}

type MoviesControllerInterface interface {
	Search(c *gin.Context)
}

func NewMoviesController(moviesService movies_service.MoviesServiceInterface) *moviesController {
	m := &moviesController{
		moviesService: moviesService,
	}

	return m
}

func (m *moviesController) Search(c *gin.Context) {
	queryParams := make(map[string]string)

	for queryKey, queryVal := range c.Request.URL.Query() {
		queryParams[queryKey] = queryVal[0]
	}

	queryParams[movies_service.QueryParam] = strings.ReplaceAll(queryParams[movies_service.QueryParam], "+", " ")

	result, searchErr := m.moviesService.SearchMovies(queryParams)
	if searchErr != nil {
		c.JSON(searchErr.Status, searchErr)

		return
	}

	c.JSON(http.StatusOK, result)
}
