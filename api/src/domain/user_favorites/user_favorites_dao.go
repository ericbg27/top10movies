package user_favorites

import (
	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

const (
	queryGetUserFavorites     = "SELECT movie_id FROM user_favorites WHERE user_id=$1"
	queryGetUserFavoritesName = "query-get-user-favorites"
)

func (u *UserFavorites) GetFavorites(userId int64) *rest_errors.RestErr {
	_, err := db.Client.Prepare(queryGetUserFavoritesName, queryGetUserFavorites)
	if err != nil {
		logger.Error("Error when trying to prepare get user favorites statement", err)
		return rest_errors.NewInternalServerError("Error when trying to get user favorites")
	}

	result, err := db.Client.Query(queryGetUserFavoritesName, userId)
	if err != nil {
		logger.Error("Error when trying to get user favorites", err)
		return rest_errors.NewInternalServerError("Error when trying to get user favorites")
	}

	for result.Next() {
		var movieId int64
		err := result.Scan(&movieId)
		if err != nil {
			logger.Error("Error when trying to get user favorites IDs", err)
			return rest_errors.NewInternalServerError("Error when trying to get user favorites")
		}

		u.MoviesIDs = append(u.MoviesIDs, movieId)
	}

	return nil
}
