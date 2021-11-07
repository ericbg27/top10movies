package app

import (
	"github.com/ericbg27/top10movies-api/src/controllers/movies"
	"github.com/ericbg27/top10movies-api/src/controllers/users"
)

func mapUrls() {
	router.POST("/login", users.Login)
	router.POST("/register", users.Create)
	router.POST("/users/:user_id", users.Update)
	router.PATCH("/users/:user_id", users.Update)
	router.DELETE("/users/:user_id", users.Delete)
	router.GET("/users/search", users.Search)

	router.GET("/users/:user_id/favorites", users.GetFavorites)
	router.POST("/users/:user_id/favorite", users.AddFavorite) // TODO: Do we put movie_id in the URL?

	router.GET("/search", movies.Search)
}
