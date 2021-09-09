package movies_service

import (
	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/ryanbradynd05/go-tmdb"
)

type moviesService struct{}

type moviesServiceInterface interface {
	SearchMovies(name, page string) (*tmdb.MovieSearchResults, *rest_errors.RestErr)
}

var (
	UsersService moviesServiceInterface = &moviesService{}

	tmdbAPI *tmdb.TMDb
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

// TODO: Add more options other than page
func (m *moviesService) SearchMovies(name, page string) (*tmdb.MovieSearchResults, *rest_errors.RestErr) {
	searchOptions := map[string]string{
		"page": page,
	}
	result, err := tmdbAPI.SearchMovie(name, searchOptions)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("Failed to search for movie")
	}

	return result, nil
}
