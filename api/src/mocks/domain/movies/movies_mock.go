package movies

import (
	movies "github.com/ericbg27/top10movies-api/src/domain/movies"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/ryanbradynd05/go-tmdb"
)

type MovieInfoMock struct {
	Movie      tmdb.Movie
	AddedMovie bool
	CanAdd     bool
	CanGet     bool
}

func (m *MovieInfoMock) AddMovie() *rest_errors.RestErr {
	if !m.CanAdd {
		return rest_errors.NewInternalServerError("Failed to add movie")
	}

	m.AddedMovie = true

	return nil
}

func (m *MovieInfoMock) GetMovie() (movies.MovieInterface, *rest_errors.RestErr) {
	if !m.CanGet {
		return nil, rest_errors.NewInternalServerError("Failed to get movie")
	}

	movie := &MovieInfoMock{
		Movie: tmdb.Movie{
			ID:    m.Movie.ID,
			Title: "Movie Test Title",
		},
	}

	return movie, nil
}
