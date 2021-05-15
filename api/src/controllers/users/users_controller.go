package users

import (
	"net/http"
	"time"

	"github.com/ericbg27/top10movies-api/src/domain/users"
	users_service "github.com/ericbg27/top10movies-api/src/services/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/gin-gonic/gin"
)

const (
	layoutISO = "2006-01-02"
)

func Login(c *gin.Context) {

}

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)

		return
	}

	user.Status = users.StatusActive
	user.DateCreated = time.Now().Format(layoutISO)

	result, saveErr := users_service.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)

		return
	}

	c.JSON(http.StatusCreated, &result)
}
