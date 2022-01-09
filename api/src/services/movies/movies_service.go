package movies_service

import (
	"github.com/ericbg27/top10movies-api/src/domain/movies"
	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/ryanbradynd05/go-tmdb"

	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
)

type moviesService struct {
	redisClient *redisdb.RedisClient
	tmdbAPI     *tmdb.TMDb
}

type MoviesServiceInterface interface {
	SearchMovies(searchOptions map[string]string) (*tmdb.MovieSearchResults, *rest_errors.RestErr)
	AddMovie(movies.MovieInterface) *rest_errors.RestErr
	GetMovieFromCache(movies.MovieInterface) (movies.MovieInterface, *rest_errors.RestErr)
	GetMovieById(int) (*tmdb.Movie, *rest_errors.RestErr)
}

const (
	QueryParam = "query"
)

func NewMoviesService(redisClient *redisdb.RedisClient, cfg *config.Config) *moviesService {
	tmdbConfig := tmdb.Config{
		APIKey:   cfg.MovieApi.ApiKey,
		Proxies:  nil,
		UseProxy: false,
	}

	tmdbAPI := tmdb.Init(tmdbConfig)

	ms := &moviesService{
		redisClient: redisClient,
		tmdbAPI:     tmdbAPI,
	}

	return ms
}

func (m *moviesService) SearchMovies(searchOptions map[string]string) (*tmdb.MovieSearchResults, *rest_errors.RestErr) {
	movieName := searchOptions[QueryParam]
	delete(searchOptions, QueryParam)

	result, err := m.tmdbAPI.SearchMovie(movieName, searchOptions)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("Failed to search for movie")
	}

	return result, nil
}

func (m *moviesService) AddMovie(movie movies.MovieInterface) *rest_errors.RestErr {
	if err := movie.AddMovie(m.redisClient); err != nil {
		return err
	}

	return nil
}

func (m *moviesService) GetMovieFromCache(movie movies.MovieInterface) (movies.MovieInterface, *rest_errors.RestErr) {
	savedMovie, err := movie.GetMovie(m.redisClient)
	if err != nil {
		return nil, err
	}

	return savedMovie, nil
}

func (m *moviesService) GetMovieById(movieId int) (*tmdb.Movie, *rest_errors.RestErr) {
	result, err := m.tmdbAPI.GetMovieInfo(movieId, nil)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("Failed to get movie information")
	}

	return result, nil
}
