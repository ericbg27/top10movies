package user_favorites

import "github.com/ericbg27/top10movies-api/src/utils/rest_errors"

type UserFavoritesInterface interface {
	GetFavorites(int64) (UserFavoritesInterface, *rest_errors.RestErr)
}

type UserFavorites struct {
	UserID    int64   `json:"user_id"`
	MoviesIDs []int64 `json:"favorite_movies"`
}
