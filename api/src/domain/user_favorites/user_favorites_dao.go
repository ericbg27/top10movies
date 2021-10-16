package user_favorites

import (
	"context"
	"fmt"
	"time"

	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	"github.com/ericbg27/top10movies-api/src/domain/movies"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

const (
	queryGetUserFavoritesIds    = "SELECT movie_id FROM user_favorites WHERE user_id=$1;"
	queryGetUserCachedFavorites = `SELECT m1.id, m1.original_title, m1.adult, m1.release_date, m1.created_at, m1.title, m1.overview FROM movies 
									AS m1 INNER JOIN (SELECT movie_id FROM user_favorites WHERE user_id=$1) AS m2 ON m1.id = m2.movie_id;`

	queryAddUserFavorite     = "INSERT INTO user_favorites VALUES ($1,$2);"
	queryAddUserFavoriteName = "query-add-user-favorite"

	releaseDateLayout = "2006-02-01"
	createdAtLayout   = "2006-02-01T15:04:05Z"
)

func (u UserFavorites) GetFavorites() (UserFavoritesInterface, map[int]bool, *rest_errors.RestErr) {
	result, err := db.Client.Query(context.Background(), queryGetUserFavoritesIds, u.UserID)
	if err != nil {
		logger.Error("Error when trying to get user favorites", err)
		return nil, nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
	}

	var userFavorites UserFavorites

	for result.Next() {
		var movieId int
		err := result.Scan(&movieId)
		if err != nil {
			logger.Error("Error when trying to get user favorites IDs", err)
			return nil, nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
		}

		userFavorites.MoviesIDs = append(userFavorites.MoviesIDs, movieId)
	}

	result, err = db.Client.Query(context.Background(), queryGetUserCachedFavorites, u.UserID)
	if err != nil { // Do we throw an error here? Maybe just log!
		logger.Error("Error when trying to get user favorites", err)
		return nil, nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
	}

	cachedIds := make(map[int]bool)

	for result.Next() {
		var cachedFavorite movies.MovieInfo
		var releaseDate time.Time
		var createdAt time.Time

		err := result.Scan(&cachedFavorite.Movie.ID, &cachedFavorite.Movie.OriginalTitle, &cachedFavorite.Movie.Adult,
			&releaseDate, &createdAt, &cachedFavorite.Movie.Title, &cachedFavorite.Movie.Overview)
		if err != nil {
			logger.Error("Error when trying to get cached user favorites", err)
			return nil, nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
		}

		cachedFavorite.Movie.ReleaseDate = releaseDate.Format(releaseDateLayout)
		cachedFavorite.CreatedAt = createdAt.Format(createdAtLayout)

		cachedIds[cachedFavorite.Movie.ID] = true
		userFavorites.MoviesData = append(userFavorites.MoviesData, cachedFavorite.Movie)
	}

	return userFavorites, cachedIds, nil
}

func (u UserFavorites) AddFavorite() *rest_errors.RestErr {
	result, err := db.Client.Exec(context.Background(), queryAddUserFavorite, u.UserID, u.MoviesIDs[0])
	if err != nil {
		logger.Error("Error when trying to prepare add user favorite statement", err)
		return rest_errors.NewBadRequestError("Error when trying to add user favorite")
	}

	logger.Info(fmt.Sprintf("Saved user in the database. Rows affected: %d", result.RowsAffected()))

	return nil
}
