package movies

import (
	"net/http"
	"strings"

	movies_service "github.com/ericbg27/top10movies-api/src/services/movies"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	queryParams := make(map[string]string)

	for queryKey, queryVal := range c.Request.URL.Query() {
		queryParams[queryKey] = queryVal[0]
	}

	queryParams[movies_service.QueryParam] = strings.ReplaceAll(queryParams[movies_service.QueryParam], "+", " ")

	result, searchErr := movies_service.UsersService.SearchMovies(queryParams)
	if searchErr != nil {
		c.JSON(searchErr.Status, searchErr)
	}

	c.JSON(http.StatusOK, result)
}
