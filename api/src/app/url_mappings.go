package app

import (
	"github.com/ericbg27/top10movies-api/src/controllers/users"
)

func mapUrls() {
	router.POST("/login", users.Login)
}
