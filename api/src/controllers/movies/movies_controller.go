package movies

import (
	"net/http"
	"strconv"
	"strings"

	movies_service "github.com/ericbg27/top10movies-api/src/services/movies"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	query := c.Query("query")
	page := c.Query("page")
	_, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, rest_errors.NewInternalServerError("Page query parameter should be an integer"))
	}

	query = strings.ReplaceAll(query, "+", " ")

	result, searchErr := movies_service.UsersService.SearchMovies(query, page)
	if searchErr != nil {
		c.JSON(searchErr.Status, searchErr)
	}

	c.JSON(http.StatusOK, result)
}
