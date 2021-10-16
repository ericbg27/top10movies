package user_favorites

const (
	QueryGetUserFavoritesIds    = "SELECT movie_id FROM user_favorites WHERE user_id=$1;"
	QueryGetUserCachedFavorites = `SELECT m1.id, m1.original_title, m1.adult, m1.release_date, m1.created_at, m1.title, m1.overview FROM movies 
									AS m1 INNER JOIN (SELECT movie_id FROM user_favorites WHERE user_id=$1) AS m2 ON m1.id = m2.movie_id;`

	QueryAddUserFavorite     = "INSERT INTO user_favorites VALUES ($1,$2);"
	QueryAddUserFavoriteName = "query-add-user-favorite"

	ReleaseDateLayout = "2006-02-01"
	CreatedAtLayout   = "2006-02-01T15:04:05Z"
)
