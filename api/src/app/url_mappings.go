package app

import (
	"github.com/ericbg27/top10movies-api/src/controllers/users"
)

func mapUrls() {
	router.POST("/login", users.Login)
	router.POST("/register", users.Create)
	router.POST("/users/:user_id", users.Update)
	router.PATCH("/users/:user_id", users.Update)
}
