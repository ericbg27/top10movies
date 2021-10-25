package user_favorites

const (
	QueryGetUserFavoritesIds = "SELECT movie_id FROM user_favorites WHERE user_id=$1;"

	QueryAddUserFavorite     = "INSERT INTO user_favorites VALUES ($1,$2);"
	QueryAddUserFavoriteName = "query-add-user-favorite"
)
