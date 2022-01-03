package user_favorites

import (
	"github.com/ericbg27/top10movies-api/src/datasources/database"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/ryanbradynd05/go-tmdb"
)

type UserFavoritesInterface interface {
	GetFavorites(database.DatabaseClient) (UserFavoritesInterface, map[int]bool, *rest_errors.RestErr)
	AddFavorite(database.DatabaseClient) *rest_errors.RestErr
}

type UserFavorites struct {
	UserID     int64        `json:"user_id"`
	MoviesIDs  []int        `json:"favorite_movies"`
	MoviesData []tmdb.Movie `json:"favorite_movies_data"`
}
