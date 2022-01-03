package app

import (
	"github.com/ericbg27/top10movies-api/src/controllers/movies"
	"github.com/ericbg27/top10movies-api/src/controllers/users"
)

func mapUrls() {
	router.POST("/login", users.UsersController.Login)
	router.POST("/register", users.UsersController.Create)
	router.POST("/users/:user_id", users.UsersController.Update)
	router.PATCH("/users/:user_id", users.UsersController.Update)
	router.DELETE("/users/:user_id", users.UsersController.Delete)
	router.GET("/users/search", users.UsersController.Search)

	router.GET("/users/:user_id/favorites", users.UsersController.GetFavorites)
	router.POST("/users/:user_id/favorite", users.UsersController.AddFavorite) // TODO: Do we put movie_id in the URL?

	router.GET("/search", movies.Search)
}
