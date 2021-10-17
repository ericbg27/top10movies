package movies

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	movies_service_mock "github.com/ericbg27/top10movies-api/src/mocks/services/movies"
	movies_service "github.com/ericbg27/top10movies-api/src/services/movies"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/ryanbradynd05/go-tmdb"
	"github.com/stretchr/testify/assert"
)

var (
	c *gin.Context
)

func PrepareTest(request []byte, method string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
		URL: &url.URL{
			RawQuery: "query=Test+Movie",
		},
	}

	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(request))

	return w
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	oldMoviesService := movies_service.MoviesService

	movies_service.MoviesService = &movies_service_mock.MoviesServiceMock{
		CanAddMovie:    true,
		CanGetMovie:    true,
		HasMovieCached: true,
		AddedMovie:     false,
		CanSearch:      true,
	}

	exitCode := m.Run()

	movies_service.MoviesService = oldMoviesService

	os.Exit(exitCode)
}

func TestSearchSuccess(t *testing.T) {
	w := PrepareTest(make([]byte, 0), "GET")

	Search(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var result tmdb.MovieSearchResults
	err := json.Unmarshal(responseData, &result)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, 1, result.Page)
	assert.EqualValues(t, 1, result.TotalPages)
	assert.EqualValues(t, 1, result.TotalResults)
	assert.EqualValues(t, 1, len(result.Results))
	assert.EqualValues(t, 1, result.Results[0].ID)
	assert.EqualValues(t, "Test Movie", result.Results[0].Title)
}

func TestSearchFail(t *testing.T) {
	w := PrepareTest(make([]byte, 0), "GET")

	movies_service.MoviesService.(*movies_service_mock.MoviesServiceMock).CanSearch = false

	Search(c)

	movies_service.MoviesService.(*movies_service_mock.MoviesServiceMock).CanSearch = true

	responseData, _ := ioutil.ReadAll(w.Body)

	var result rest_errors.RestErr
	err := json.Unmarshal(responseData, &result)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, http.StatusInternalServerError, result.Status)
	assert.EqualValues(t, "Failed to search for movies", result.Message)
	assert.EqualValues(t, "internal_server_error", result.Err)
}
