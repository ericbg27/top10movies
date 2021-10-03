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

func (m Movie) AddMovie() *rest_errors.RestErr {
	/*_, err := db.Client.Prepare(queryAddMovieName, queryAddMovie)
	if err != nil {
		logger.Error("Error when trying to prepare add movie statement", err)
		return rest_errors.NewInternalServerError("Error when trying to add movie")
	}*/

	var releaseDate time.Time
	var createdAt time.Time

	releaseDate, _ = time.Parse(releaseDateLayout, m.ReleaseDate)
	createdAt = time.Now()

	_, err := db.Client.Exec(context.Background(), queryAddMovie, m.ID, m.OriginalTitle, m.Adult, releaseDate, createdAt, m.Title, m.Overview)
	if err != nil {
		logger.Error("Error when trying to add movie", err)
		return rest_errors.NewInternalServerError("Error when trying to add movie")
	}

	return nil
}

func (m Movie) GetMovie() (MovieInterface, *rest_errors.RestErr) {
	/*_, err := db.Client.Prepare(queryGetMovieName, queryGetMovie)
	if err != nil {
		logger.Error("Error when trying to prepare get movie statement", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
	}*/

	var savedMovie Movie

	result := db.Client.QueryRow(context.Background(), queryGetMovie, m.ID)

	var releaseDate time.Time
	var createdAt time.Time

	err := result.Scan(&savedMovie.ID, &savedMovie.OriginalTitle, &savedMovie.Adult, &releaseDate, &createdAt, &savedMovie.Title, &savedMovie.Overview)

	savedMovie.ReleaseDate = releaseDate.Format(releaseDateLayout)
	savedMovie.CreatedAt = createdAt.Format(createdAtLayout)

	if err != nil && err != pgx.ErrNoRows {
		logger.Error("Error when trying to get movie", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
	} else if err == pgx.ErrNoRows {
		savedMovie.ID = -1
	}

	return savedMovie, nil
}
