package movies

import (
	"context"
	"time"

	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/jackc/pgx/v4"
)

const (
	queryAddMovie = "INSERT INTO movies VALUES ($1,$2,$3,$4,$5,$6,$7);"
	queryGetMovie = "SELECT * FROM movies WHERE id=$1;"

	releaseDateLayout = "2006-02-01"
	createdAtLayout   = "2006-02-01T15:04:05Z"
)

// TODO: Add more informtion to the movies table in database
func (m MovieInfo) AddMovie() *rest_errors.RestErr {
	var releaseDate time.Time
	var createdAt time.Time

	releaseDate, _ = time.Parse(releaseDateLayout, m.Movie.ReleaseDate)
	createdAt = time.Now()

	_, err := db.Client.Exec(context.Background(), queryAddMovie, m.Movie.ID, m.Movie.OriginalTitle, m.Movie.Adult, releaseDate, createdAt, m.Movie.Title, m.Movie.Overview)
	if err != nil {
		logger.Error("Error when trying to add movie", err)
		return rest_errors.NewInternalServerError("Error when trying to add movie")
	}

	return nil
}

func (m MovieInfo) GetMovie() (MovieInterface, *rest_errors.RestErr) {
	var savedMovie MovieInfo

	result := db.Client.QueryRow(context.Background(), queryGetMovie, m.Movie.ID)

	var releaseDate time.Time
	var createdAt time.Time

	err := result.Scan(&savedMovie.Movie.ID, &savedMovie.Movie.OriginalTitle, &savedMovie.Movie.Adult, &releaseDate, &createdAt, &savedMovie.Movie.Title, &savedMovie.Movie.Overview)

	savedMovie.Movie.ReleaseDate = releaseDate.Format(releaseDateLayout)
	savedMovie.CreatedAt = createdAt.Format(createdAtLayout)

	if err != nil && err != pgx.ErrNoRows {
		logger.Error("Error when trying to get movie", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
	} else if err == pgx.ErrNoRows {
		savedMovie.Movie.ID = -1
	}

	return savedMovie, nil
}
