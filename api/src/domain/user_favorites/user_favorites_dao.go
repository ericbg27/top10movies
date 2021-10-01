package user_favorites

import (
	"fmt"

	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

const (
	queryGetUserFavorites     = "SELECT movie_id FROM user_favorites WHERE user_id=$1;"
	queryGetUserFavoritesName = "query-get-user-favorites"

	queryAddUserFavorite     = "INSERT INTO user_favorites VALUES ($1,$2);"
	queryAddUserFavoriteName = "query-add-user-favorite"
)

func (u UserFavorites) GetFavorites() (UserFavoritesInterface, *rest_errors.RestErr) {
	_, err := db.Client.Prepare(queryGetUserFavoritesName, queryGetUserFavorites)
	if err != nil {
		logger.Error("Error when trying to prepare get user favorites statement", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
	}

	result, err := db.Client.Query(queryGetUserFavoritesName, u.UserID)
	if err != nil {
		logger.Error("Error when trying to get user favorites", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
	}

	var userFavorites UserFavorites

	for result.Next() {
		var movieId int64
		err := result.Scan(&movieId)
		if err != nil {
			logger.Error("Error when trying to get user favorites IDs", err)
			return nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
		}

		userFavorites.MoviesIDs = append(userFavorites.MoviesIDs, movieId)
	}

	return userFavorites, nil
}

func (u UserFavorites) AddFavorite() *rest_errors.RestErr {
	_, err := db.Client.Prepare(queryAddUserFavoriteName, queryAddUserFavorite)
	if err != nil {
		logger.Error("Error when trying prepare add user favorite statement", err)
		return rest_errors.NewInternalServerError("Error when trying to add user favorite")
	}

	result, err := db.Client.Exec(queryAddUserFavoriteName, u.UserID, u.MoviesIDs[0])
	if err != nil {
		logger.Error("Error when trying to prepare add user favorite statement", err)
		return rest_errors.NewBadRequestError("Error when trying to add user favorite")
	}

	logger.Info(fmt.Sprintf("Saved user in the database. Rows affected: %d", result.RowsAffected()))

	return nil
}
