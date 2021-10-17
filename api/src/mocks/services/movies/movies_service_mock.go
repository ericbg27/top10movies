package movies_service

import (
	"github.com/ericbg27/top10movies-api/src/domain/movies"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/ryanbradynd05/go-tmdb"
)

type MoviesServiceMock struct {
	CanAddMovie    bool
	CanGetMovie    bool
	HasMovieCached bool
	CanSearch      bool
	AddedMovie     bool
}

func (m *MoviesServiceMock) SearchMovies(searchOptions map[string]string) (*tmdb.MovieSearchResults, *rest_errors.RestErr) {
	query := searchOptions["query"]

	if !m.CanSearch {
		return nil, rest_errors.NewInternalServerError("Failed to search for movies")
	}

	var result tmdb.MovieSearchResults
	result.Page = 1
	result.TotalPages = 1
	result.Results = append(result.Results, tmdb.MovieShort{
		ID:    1,
		Title: query,
	})
	result.TotalResults = 1

	return &result, nil
}

func (m *MoviesServiceMock) AddMovie(movie movies.MovieInterface) *rest_errors.RestErr {
	if !m.CanAddMovie {
		return rest_errors.NewInternalServerError("Error when trying to add movie")
	}

	m.AddedMovie = true

	return nil
}

func (m *MoviesServiceMock) GetMovieFromCache(movie movies.MovieInterface) (movies.MovieInterface, *rest_errors.RestErr) {
	mov := movie.(movies.MovieInfo)

	if !m.CanGetMovie {
		return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
	}

	if m.HasMovieCached {
		mov.Movie = tmdb.Movie{
			Title: "Example Movie Title",
			ID:    1,
		}
		mov.CreatedAt = "01-02-2006"
	} else {
		mov.Movie = tmdb.Movie{
			Title: "",
			ID:    -1,
		}
		mov.CreatedAt = ""
	}

	return mov, nil
}

func (m *MoviesServiceMock) GetMovieById(movieId int) (*tmdb.Movie, *rest_errors.RestErr) {
	if !m.CanGetMovie {
		return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
	}

	movieInfo := &tmdb.Movie{
		ID: movieId,
	}

	return movieInfo, nil
}
