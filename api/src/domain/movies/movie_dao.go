package movies

import (
	"context"
	"time"

	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	movies_queries "github.com/ericbg27/top10movies-api/src/queries/movies"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/jackc/pgx/v4"
)

// TODO: Add more informtion to the movies table in database
func (m MovieInfo) AddMovie() *rest_errors.RestErr {
	var releaseDate time.Time
	var createdAt time.Time

	releaseDate, _ = time.Parse(movies_queries.ReleaseDateLayout, m.Movie.ReleaseDate)
	createdAt = time.Now()

	_, err := db.Client.Exec(context.Background(), movies_queries.QueryAddMovie, m.Movie.ID, m.Movie.OriginalTitle, m.Movie.Adult, releaseDate, createdAt, m.Movie.Title, m.Movie.Overview)
	if err != nil {
		logger.Error("Error when trying to add movie", err)
		return rest_errors.NewInternalServerError("Error when trying to add movie")
	}

	return nil
}

func (m MovieInfo) GetMovie() (MovieInterface, *rest_errors.RestErr) {
	var savedMovie MovieInfo

	result := db.Client.QueryRow(context.Background(), movies_queries.QueryGetMovie, m.Movie.ID)

	var releaseDate time.Time
	var createdAt time.Time

	err := result.Scan(&savedMovie.Movie.ID, &savedMovie.Movie.OriginalTitle, &savedMovie.Movie.Adult, &releaseDate, &createdAt, &savedMovie.Movie.Title, &savedMovie.Movie.Overview)

	savedMovie.Movie.ReleaseDate = releaseDate.Format(movies_queries.ReleaseDateLayout)
	savedMovie.CreatedAt = createdAt.Format(movies_queries.CreatedAtLayout)

	if err != nil && err != pgx.ErrNoRows {
		logger.Error("Error when trying to get movie", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
	} else if err == pgx.ErrNoRows {
		savedMovie.Movie.ID = -1
	}

	return savedMovie, nil
}
