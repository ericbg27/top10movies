package users

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ericbg27/top10movies-api/src/domain/movies"
	"github.com/ericbg27/top10movies-api/src/domain/user_favorites"
	"github.com/ericbg27/top10movies-api/src/domain/users"
	movies_service "github.com/ericbg27/top10movies-api/src/services/movies"
	users_service "github.com/ericbg27/top10movies-api/src/services/users"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	layoutISO = "2006-01-02"
)

func getID(userIDParam string) (int64, *rest_errors.RestErr) {
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

	// TODO: Create a token and send to the user
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
	userID, IdErr := getID(c.Param("user_id"))
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

func Delete(c *gin.Context) {
	userID, IdErr := getID(c.Param("user_id"))
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

	deleteErr := users_service.UsersService.DeleteUser(user)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)

		return
	}

	c.Status(http.StatusOK)
}

func GetUserFavorites(c *gin.Context) {
	userID, IdErr := getID(c.Param("user_id"))
	if IdErr != nil {
		c.JSON(IdErr.Status, IdErr)

		return
	}

	var usrFav user_favorites.UserFavorites
	usrFav.UserID = userID

	userFavorites, getErr := users_service.UsersService.GetUserFavorites(usrFav)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)

		return
	}

	c.JSON(http.StatusOK, userFavorites)
}

func AddUserFavorite(c *gin.Context) {
	userID, userIdErr := getID(c.Param("user_id"))
	if userIdErr != nil {
		c.JSON(userIdErr.Status, userIdErr)

		return
	}

	var movie movies.MovieInfo

	if err := c.ShouldBindJSON(&movie); err != nil {
		restErr := rest_errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)

		return
	}

	movieCacheResult, err := movies_service.MoviesService.GetMovie(movie)
	if err != nil {
		c.JSON(err.Status, err)

		return
	}

	movieCache := movieCacheResult.(movies.MovieInfo)
	if movieCache.Movie.ID == -1 { // Movie is not cached
		addErr := movies_service.MoviesService.AddMovie(movie)
		if addErr != nil {
			c.JSON(addErr.Status, addErr)

			return
		}
	}

	var userFavorite user_favorites.UserFavorites
	userFavorite.UserID = userID
	userFavorite.MoviesIDs = append(userFavorite.MoviesIDs, movie.Movie.ID)

	addErr := users_service.UsersService.AddUserFavorite(userFavorite)
	if addErr != nil {
		c.JSON(addErr.Status, addErr)

		return
	}

	c.Status(http.StatusOK)
}
