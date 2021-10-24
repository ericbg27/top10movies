package user_favorites

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
	"github.com/ericbg27/top10movies-api/src/domain/movies"
	user_favorites_queries "github.com/ericbg27/top10movies-api/src/queries/user_favorites"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

func (u UserFavorites) GetFavorites() (UserFavoritesInterface, map[int]bool, *rest_errors.RestErr) {
	result, err := db.Client.Query(context.Background(), user_favorites_queries.QueryGetUserFavoritesIds, u.UserID)
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

	cachedIds := make(map[int]bool)

	for _, movieId := range userFavorites.MoviesIDs {
		var movieRedisKey strings.Builder
		movieRedisKey.WriteString("movie:")
		movieRedisKey.WriteString(strconv.Itoa(movieId))

		redisResult, redisErr := redisdb.Client.Get(movieRedisKey.String()).Result()
		if redisErr != nil && redisErr != redisdb.RedisNil { // Do we throw an error here? Maybe just log!
			logger.Error("Error when trying to get user favorites", redisErr)
			return nil, nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
		}

		if redisErr != redisdb.RedisNil {
			var cachedFavorite movies.MovieInfo
			err = json.Unmarshal([]byte(redisResult), &cachedFavorite)
			if err != nil {
				logger.Error("Error when trying to get user favorites", err)
				return nil, nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
			}

			cachedIds[movieId] = true
			userFavorites.MoviesData = append(userFavorites.MoviesData, cachedFavorite.Movie)
		}
	}

	return userFavorites, cachedIds, nil
}

func (u UserFavorites) AddFavorite() *rest_errors.RestErr {
	result, err := db.Client.Exec(context.Background(), user_favorites_queries.QueryAddUserFavorite, u.UserID, u.MoviesIDs[0])
	if err != nil {
		logger.Error("Error when trying to prepare add user favorite statement", err)
		return rest_errors.NewBadRequestError("Error when trying to add user favorite")
	}

	logger.Info(fmt.Sprintf("Saved user in the database. Rows affected: %d", result.RowsAffected()))

	return nil
}
