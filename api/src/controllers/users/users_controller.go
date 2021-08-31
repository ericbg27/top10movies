package users

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ericbg27/top10movies-api/src/domain/users"
	users_service "github.com/ericbg27/top10movies-api/src/services/users"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	layoutISO = "2006-01-02"
)

func getUserID(userIDParam string) (int64, *rest_errors.RestErr) {
	userID, userErr := strconv.ParseInt(userIDParam, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("User ID should be a number")
	}

	return userID, nil
}

func Login(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)

		return
	}

	result, getErr := users_service.UsersService.GetUser(user)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)

		return
	}

	savedUser := result.(users.User)

	err := bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password))
	if err != nil {
		passwordErr := rest_errors.NewBadRequestError("Wrong password")
		c.JSON(passwordErr.Status, passwordErr)

		return
	}

	// Here we create a token and send to the user
	c.JSON(http.StatusOK, nil)
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

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Unable to hash password", err)
		hashErr := rest_errors.NewBadRequestError("Unable to hash password")
		c.JSON(hashErr.Status, hashErr)

		return
	}

	user.Password = string(hashedPass)

	result, saveErr := users_service.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)

		return
	}

	newUser := result.(users.User)

	newUser.Password = ""

	c.JSON(http.StatusCreated, newUser)
}

func Update(c *gin.Context) {
	userID, IdErr := getUserID(c.Param("user_id"))
	if IdErr != nil {
		c.JSON(IdErr.Status, IdErr)

		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)

		return
	}

	user.ID = userID

	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := users_service.UsersService.UpdateUser(user, isPartial)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)

		return
	}

	updatedUser := result.(users.User)
	updatedUser.Password = ""

	c.JSON(http.StatusOK, updatedUser)
}
