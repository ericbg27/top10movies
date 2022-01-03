package users

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericbg27/top10movies-api/src/domain/movies"
	"github.com/ericbg27/top10movies-api/src/domain/user_favorites"
	"github.com/ericbg27/top10movies-api/src/domain/users"
	movies_service "github.com/ericbg27/top10movies-api/src/services/movies"
	users_service "github.com/ericbg27/top10movies-api/src/services/users"
	"github.com/ericbg27/top10movies-api/src/utils/authorization"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	layoutISO = "2006-01-02"
)

type usersController struct{}

type UsersControllerInterface interface {
	Login(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetFavorites(c *gin.Context)
	AddFavorite(c *gin.Context)
	Search(c *gin.Context)
}

var (
	UsersController UsersControllerInterface = &usersController{}
)

func getID(userIDParam string) (int64, *rest_errors.RestErr) {
	userID, userErr := strconv.ParseInt(userIDParam, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("User ID should be a number")
	}

	return userID, nil
}

func (u *usersController) Login(c *gin.Context) {
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

	token, err := authorization.AuthManager.CreateToken(savedUser.ID)
	if err != nil {
		tokenErr := rest_errors.NewInternalServerError("Could not generate jwt access token")
		c.JSON(tokenErr.Status, tokenErr)

		return
	}

	tokensInfo := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	c.JSON(http.StatusOK, tokensInfo)
}

func (u *usersController) Create(c *gin.Context) {
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

func (u *usersController) Update(c *gin.Context) {
	bearToken := c.Request.Header.Get("Authorization")

	userID, err := authorization.AuthManager.FetchAuth(bearToken)
	if err != nil {
		authErr := rest_errors.NewUnauthorizedError("Invalid JWT token")
		c.JSON(authErr.Status, authErr)

		return
	}

	requestUserID, IdErr := getID(c.Param("user_id"))
	if IdErr != nil {
		c.JSON(IdErr.Status, IdErr)

		return
	}

	if requestUserID != int64(userID) {
		wrongIdErr := rest_errors.NewUnauthorizedError("User ID in the request does not match token user ID")
		c.JSON(wrongIdErr.Status, wrongIdErr)

		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)

		return
	}

	user.ID = int64(userID)

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

func (u *usersController) Delete(c *gin.Context) {
	bearToken := c.Request.Header.Get("Authorization")

	userID, err := authorization.AuthManager.FetchAuth(bearToken)
	if err != nil {
		authErr := rest_errors.NewUnauthorizedError("Invalid JWT token")
		c.JSON(authErr.Status, authErr)

		return
	}

	requestUserID, IdErr := getID(c.Param("user_id"))
	if IdErr != nil {
		c.JSON(IdErr.Status, IdErr)

		return
	}

	if requestUserID != int64(userID) {
		wrongIdErr := rest_errors.NewUnauthorizedError("User ID in the request does not match token user ID")
		c.JSON(wrongIdErr.Status, wrongIdErr)

		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)

		return
	}

	user.ID = int64(userID)

	deleteErr := users_service.UsersService.DeleteUser(user)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)

		return
	}

	c.Status(http.StatusOK)
}

func (u *usersController) GetFavorites(c *gin.Context) {
	userID, IdErr := getID(c.Param("user_id"))
	if IdErr != nil {
		c.JSON(IdErr.Status, IdErr)

		return
	}

	var usrFav user_favorites.UserFavorites
	usrFav.UserID = userID

	userFavorites, cachedFavorites, getErr := users_service.UsersService.GetUserFavorites(usrFav)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)

		return
	}

	usrFav.MoviesData = append(usrFav.MoviesData, userFavorites.(user_favorites.UserFavorites).MoviesData...)

	for _, movieId := range userFavorites.(user_favorites.UserFavorites).MoviesIDs {
		if _, cached := cachedFavorites[movieId]; !cached {
			var movie movies.MovieInfo
			movie.Movie.ID = movieId

			movieResult, err := movies_service.MoviesService.GetMovieById(movie.Movie.ID)
			if err != nil { // TODO: Do we return an error if one of the favorites is not found?
				c.JSON(err.Status, err)

				return
			}

			movie.Movie = *movieResult
			addErr := movies_service.MoviesService.AddMovie(movie)
			if addErr != nil { // TODO: Do we return an error if we fail to save in cache? Maybe just log!
				c.JSON(addErr.Status, addErr)

				return
			}

			usrFav.MoviesData = append(usrFav.MoviesData, *movieResult)
		}
	}

	c.JSON(http.StatusOK, usrFav)
}

func (u *usersController) AddFavorite(c *gin.Context) {
	bearToken := c.Request.Header.Get("Authorization")

	userID, err := authorization.AuthManager.FetchAuth(bearToken)
	if err != nil {
		authErr := rest_errors.NewUnauthorizedError("Invalid JWT token")
		c.JSON(authErr.Status, authErr)

		return
	}

	requestUserID, IdErr := getID(c.Param("user_id"))
	if IdErr != nil {
		c.JSON(IdErr.Status, IdErr)

		return
	}

	if requestUserID != int64(userID) {
		wrongIdErr := rest_errors.NewUnauthorizedError("User ID in the request does not match token user ID")
		c.JSON(wrongIdErr.Status, wrongIdErr)

		return
	}

	var movie movies.MovieInfo

	if err := c.ShouldBindJSON(&movie); err != nil {
		restErr := rest_errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)

		return
	}

	movieCacheResult, cacheErr := movies_service.MoviesService.GetMovieFromCache(movie)
	if cacheErr != nil {
		c.JSON(cacheErr.Status, cacheErr)

		return
	}

	movieCache := movieCacheResult.(movies.MovieInfo)
	if movieCache.Movie.ID == -1 { // Movie is not cached
		addErr := movies_service.MoviesService.AddMovie(movie)
		if addErr != nil { // TODO: Do we return an error if we fail to save in cache? Maybe just log!
			c.JSON(addErr.Status, addErr)

			return
		}
	}

	var userFavorite user_favorites.UserFavorites
	userFavorite.UserID = int64(userID)
	userFavorite.MoviesIDs = append(userFavorite.MoviesIDs, movie.Movie.ID)

	addErr := users_service.UsersService.AddUserFavorite(userFavorite)
	if addErr != nil {
		c.JSON(addErr.Status, addErr)

		return
	}

	c.Status(http.StatusOK)
}

func (u *usersController) Search(c *gin.Context) {
	queryParams := make(map[string]string)

	for queryKey, queryVal := range c.Request.URL.Query() {
		queryParams[queryKey] = queryVal[0]
	}

	queryParams[users_service.QueryParam] = strings.ReplaceAll(queryParams[users_service.QueryParam], "+", " ")

	queryArray := strings.Split(queryParams[users_service.QueryParam], " ")

	var userToSearch users.User
	userToSearch.FirstName = queryArray[0]
	userToSearch.LastName = strings.Join(queryArray[1:], " ")

	foundUsers, searchErr := users_service.UsersService.SearchUser(userToSearch)
	if searchErr != nil {
		c.JSON(searchErr.Status, searchErr)

		return
	}

	c.JSON(http.StatusOK, foundUsers)
}
