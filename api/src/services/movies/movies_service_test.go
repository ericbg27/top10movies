package movies_service

import (
	"net/http"
	"testing"

	movies_mock "github.com/ericbg27/top10movies-api/src/mocks/domain/movies"
	"github.com/ryanbradynd05/go-tmdb"
	"github.com/stretchr/testify/assert"
)

func TestAddMovieSuccess(t *testing.T) {
	movieToAdd := movies_mock.MovieInfoMock{
		CanAdd:     true,
		AddedMovie: false,
	}

	addErr := MoviesService.AddMovie(&movieToAdd)

	assert.Nil(t, addErr)
	assert.EqualValues(t, true, movieToAdd.AddedMovie)
}

func TestAddMovieFailure(t *testing.T) {
	movieToAdd := movies_mock.MovieInfoMock{
		CanAdd:     false,
		AddedMovie: false,
	}

	addErr := MoviesService.AddMovie(&movieToAdd)

	assert.NotNil(t, addErr)
	assert.EqualValues(t, http.StatusInternalServerError, addErr.Status)
	assert.EqualValues(t, "Failed to add movie", addErr.Message)
	assert.EqualValues(t, "internal_server_error", addErr.Err)
}

func TestGetMovieFromCacheSuccess(t *testing.T) {
	movieToGet := movies_mock.MovieInfoMock{
		CanGet: true,
		Movie: tmdb.Movie{
			ID: 1,
		},
	}

	result, getErr := MoviesService.GetMovieFromCache(&movieToGet)

	movie := result.(*movies_mock.MovieInfoMock)

	assert.Nil(t, getErr)
	assert.NotNil(t, movie)
	assert.EqualValues(t, 1, movie.Movie.ID)
	assert.EqualValues(t, "Movie Test Title", movie.Movie.Title)
}
