package movies_service

import (
	"github.com/ericbg27/top10movies-api/src/domain/movies"
	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/ryanbradynd05/go-tmdb"
)

type moviesService struct{}

type moviesServiceInterface interface {
	SearchMovies(searchOptions map[string]string) (*tmdb.MovieSearchResults, *rest_errors.RestErr)
	AddMovie(movies.MovieInterface) *rest_errors.RestErr
	GetMovie(movies.MovieInterface) (movies.MovieInterface, *rest_errors.RestErr)
}

var (
	MoviesService moviesServiceInterface = &moviesService{}

	tmdbAPI *tmdb.TMDb
)

const (
	QueryParam = "query"
)

func init() {
	cfg := config.GetConfig()

	tmdbConfig := tmdb.Config{
		APIKey:   cfg.MovieApi.ApiKey,
		Proxies:  nil,
		UseProxy: false,
	}

	tmdbAPI = tmdb.Init(tmdbConfig)
}

func (m *moviesService) SearchMovies(searchOptions map[string]string) (*tmdb.MovieSearchResults, *rest_errors.RestErr) {
	movieName := searchOptions[QueryParam]
	delete(searchOptions, QueryParam)

	result, err := tmdbAPI.SearchMovie(movieName, searchOptions)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("Failed to search for movie")
	}

	return result, nil
}

func (m *moviesService) AddMovie(movie movies.MovieInterface) *rest_errors.RestErr {
	if err := movie.AddMovie(); err != nil {
		return err
	}

	return nil
}

func (m *moviesService) GetMovie(movie movies.MovieInterface) (movies.MovieInterface, *rest_errors.RestErr) {
	savedMovie, err := movie.GetMovie()
	if err != nil {
		return nil, err
	}

	return savedMovie, nil
}
