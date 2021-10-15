package users

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ericbg27/top10movies-api/src/domain/movies"
	"github.com/ericbg27/top10movies-api/src/domain/user_favorites"
	"github.com/ericbg27/top10movies-api/src/domain/users"
	movies_service "github.com/ericbg27/top10movies-api/src/services/movies"
	users_service "github.com/ericbg27/top10movies-api/src/services/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/ryanbradynd05/go-tmdb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	c        *gin.Context
	mockDb   map[string]string
	mockDbID map[int64]users.User
	now      string
)

type usersServiceMock struct {
	canGetFavorites bool
	canAddFavorite  bool
}

func (u *usersServiceMock) CreateUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	usr := user.(users.User)
	if _, ok := mockDb[usr.Email]; ok {
		return nil, rest_errors.NewInternalServerError("Error when trying to save user")
	}

	usr.DateCreated = now
	usr.ID = 2

	return usr, nil
}

func (u *usersServiceMock) GetUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	usr := user.(users.User)
	if savedPassword, ok := mockDb[usr.Email]; ok {
		savedUser := users.User{
			Email:    usr.Email,
			Password: savedPassword,
		}

		return savedUser, nil
	}

	return nil, rest_errors.NewNotFoundError("User not found")
}

func (u *usersServiceMock) UpdateUser(user users.UserInterface, isPartial bool) (users.UserInterface, *rest_errors.RestErr) {
	newUser := user.(users.User)

	currentUser, ok := mockDbID[newUser.ID]
	if !ok {
		return nil, rest_errors.NewInternalServerError("Error when trying to update user")
	}

	if isPartial {
		if newUser.FirstName == "" {
			newUser.FirstName = currentUser.FirstName
		}
		if newUser.LastName == "" {
			newUser.LastName = currentUser.LastName
		}
		if newUser.Email == "" {
			newUser.Email = currentUser.Email
		}
	}

	return newUser, nil
}

func (u *usersServiceMock) DeleteUser(user users.UserInterface) *rest_errors.RestErr {
	usr := user.(users.User)

	_, ok := mockDbID[usr.ID]
	if !ok {
		return rest_errors.NewInternalServerError("Error when trying to delete user")
	}

	return nil
}

func (u *usersServiceMock) GetUserFavorites(userFavs user_favorites.UserFavoritesInterface) (user_favorites.UserFavoritesInterface, *rest_errors.RestErr) {
	userFavorites := userFavs.(user_favorites.UserFavorites)

	if !u.canGetFavorites {
		return nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
	}

	userFavorites.MoviesIDs = append(userFavorites.MoviesIDs, 1)

	return userFavorites, nil
}

func (u *usersServiceMock) AddUserFavorite(userFavs user_favorites.UserFavoritesInterface) *rest_errors.RestErr {
	if !u.canAddFavorite {
		return rest_errors.NewInternalServerError("Error when trying to add user favorite")
	}

	return nil
}

type moviesServiceMock struct {
	canAddMovie    bool
	canGetMovie    bool
	hasMovieCached bool
	addedMovie     bool
}

func (m *moviesServiceMock) SearchMovies(searchOptions map[string]string) (*tmdb.MovieSearchResults, *rest_errors.RestErr) {
	return nil, nil
}

func (m *moviesServiceMock) AddMovie(movie movies.MovieInterface) *rest_errors.RestErr {
	if !m.canAddMovie {
		return rest_errors.NewInternalServerError("Error when trying to add movie")
	}

	m.addedMovie = true

	return nil
}

func (m *moviesServiceMock) GetMovieFromCache(movie movies.MovieInterface) (movies.MovieInterface, *rest_errors.RestErr) {
	mov := movie.(movies.MovieInfo)

	if !m.canGetMovie {
		return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
	}

	if m.hasMovieCached {
		mov.Movie = tmdb.Movie{
			Title: "Example Movie Title",
			ID:    1,
		}
		mov.CreatedAt = "01-02-2006"
	} else {
		mov.Movie = tmdb.Movie{
			Title: "",
			ID:    -1,
		}
		mov.CreatedAt = ""
	}

	return mov, nil
}

func (*moviesServiceMock) GetMovieById(movieId int) (*tmdb.Movie, *rest_errors.RestErr) {
	// TODO
	return nil, nil
}

func PrepareTest(request []byte, method string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(request))

	return w
}

func TestMain(m *testing.M) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	mockDb = map[string]string{
		"johndoe@gmail.com": string(hashedPass),
	}
	mockDbID = map[int64]users.User{
		1: {
			ID:          1,
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "johndoe@gmail.com",
			DateCreated: "",
			Status:      "",
			Password:    "1234",
		},
	}

	gin.SetMode(gin.TestMode)

	now = time.Now().Format(layoutISO)

	users_service.UsersService = &usersServiceMock{
		canGetFavorites: true,
		canAddFavorite:  true,
	}

	movies_service.MoviesService = &moviesServiceMock{
		canAddMovie:    true,
		canGetMovie:    true,
		hasMovieCached: true,
		addedMovie:     false,
	}

	os.Exit(m.Run())
}

func TestLoginSuccess(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			Email:    "johndoe@gmail.com",
			Password: "123456",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	Login(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, "null", string(responseData))
}

func TestLoginWrongPassword(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			Email:    "johndoe@gmail.com",
			Password: "12345",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	Login(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, "Wrong password", receivedResponse.Message)
	assert.EqualValues(t, http.StatusBadRequest, receivedResponse.Status)
	assert.EqualValues(t, "bad_request", receivedResponse.Err)
}

func TestLoginInvalidJSON(t *testing.T) {
	exampleJsonReq, err := json.Marshal(`{"invalid_key": "true"}`)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	Login(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, receivedResponse.Message, "Invalid JSON body")
	assert.EqualValues(t, receivedResponse.Status, http.StatusBadRequest)
	assert.EqualValues(t, receivedResponse.Err, "bad_request")
}

func TestLoginUserNotFound(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			Email:    "nonregisteredemail@gmail.com",
			Password: "1234",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	Login(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusNotFound, w.Code)
	assert.EqualValues(t, "User not found", receivedResponse.Message)
	assert.EqualValues(t, http.StatusNotFound, receivedResponse.Status)
	assert.EqualValues(t, "not_found", receivedResponse.Err)
}

func TestCreateSuccess(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			Email:     "johndoe2@gmail.com",
			Password:  "1234",
			FirstName: "John",
			LastName:  "Doe",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	Create(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse users.User
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, w.Code)
	assert.EqualValues(t, 2, receivedResponse.ID)
	assert.EqualValues(t, users.StatusActive, receivedResponse.Status)
	assert.EqualValues(t, "John", receivedResponse.FirstName)
	assert.EqualValues(t, "Doe", receivedResponse.LastName)
	assert.EqualValues(t, "johndoe2@gmail.com", receivedResponse.Email)
	assert.EqualValues(t, "", receivedResponse.Password)
	assert.EqualValues(t, now, receivedResponse.DateCreated)
}

func TestCreateInvalidJSON(t *testing.T) {
	exampleJsonReq, err := json.Marshal(`{"invalid_key": "true"}`)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	Create(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, receivedResponse.Message, "Invalid JSON body")
	assert.EqualValues(t, receivedResponse.Status, http.StatusBadRequest)
	assert.EqualValues(t, receivedResponse.Err, "bad_request")
}

func TestCreateSaveError(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			Email:    "johndoe@gmail.com",
			Password: "1234",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	Create(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, "Error when trying to save user", receivedResponse.Message)
	assert.EqualValues(t, http.StatusInternalServerError, receivedResponse.Status)
	assert.EqualValues(t, receivedResponse.Err, "internal_server_error")
}

func TestUpdateSuccess(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			Email:     "johndoe2@gmail.com",
			FirstName: "Johnn",
			LastName:  "Doee",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	Update(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse users.User
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, 1, receivedResponse.ID)
	assert.EqualValues(t, "Johnn", receivedResponse.FirstName)
	assert.EqualValues(t, "Doee", receivedResponse.LastName)
	assert.EqualValues(t, "johndoe2@gmail.com", receivedResponse.Email)
	assert.EqualValues(t, "", receivedResponse.Password)
}

func TestUpdatePartialSuccess(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			FirstName: "Johnn",
			LastName:  "Doee",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "PATCH")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	Update(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse users.User
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, 1, receivedResponse.ID)
	assert.EqualValues(t, "Johnn", receivedResponse.FirstName)
	assert.EqualValues(t, "Doee", receivedResponse.LastName)
	assert.EqualValues(t, "johndoe@gmail.com", receivedResponse.Email)
	assert.EqualValues(t, "", receivedResponse.Password)
}

func TestUpdateInvalidUserID(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			FirstName: "Johnn",
			LastName:  "Doee",
			Email:     "johndoe2@gmail.com",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	Update(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, "User ID should be a number", receivedResponse.Message)
	assert.EqualValues(t, "bad_request", receivedResponse.Err)
	assert.EqualValues(t, http.StatusBadRequest, receivedResponse.Status)
}

func TestUpdateInvalidJSONBody(t *testing.T) {
	exampleJsonReq, err := json.Marshal(`{"invalid_key": "true"}`)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "2"})

	Update(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, receivedResponse.Message, "Invalid JSON body")
	assert.EqualValues(t, receivedResponse.Status, http.StatusBadRequest)
	assert.EqualValues(t, receivedResponse.Err, "bad_request")
}

func TestUpdateSaveError(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			ID:        2,
			FirstName: "Johnn",
			LastName:  "Doee",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "2"})

	Update(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, "Error when trying to update user", receivedResponse.Message)
	assert.EqualValues(t, "internal_server_error", receivedResponse.Err)
	assert.EqualValues(t, http.StatusInternalServerError, receivedResponse.Status)
}

func TestDeleteSuccess(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			ID:        1,
			FirstName: "Johnn",
			LastName:  "Doee",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "DELETE")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	Delete(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	receivedResponse := string(responseData[:])

	assert.EqualValues(t, "", receivedResponse)
}

func TestDeleteInvalidUserID(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			FirstName: "Johnn",
			LastName:  "Doee",
			Email:     "johndoe2@gmail.com",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "Delete")

	Delete(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, "User ID should be a number", receivedResponse.Message)
	assert.EqualValues(t, "bad_request", receivedResponse.Err)
	assert.EqualValues(t, http.StatusBadRequest, receivedResponse.Status)
}

func TestDeleteInvalidJSONBody(t *testing.T) {
	exampleJsonReq, err := json.Marshal(`{"invalid_key": "true"}`)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "DELETE")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	Delete(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, receivedResponse.Message, "Invalid JSON body")
	assert.EqualValues(t, receivedResponse.Status, http.StatusBadRequest)
	assert.EqualValues(t, receivedResponse.Err, "bad_request")
}

func TestDeleteDeleteError(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		users.User{
			FirstName: "Johnn",
			LastName:  "Doee",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "DELETE")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "2"})

	Delete(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, "Error when trying to delete user", receivedResponse.Message)
	assert.EqualValues(t, http.StatusInternalServerError, receivedResponse.Status)
	assert.EqualValues(t, "internal_server_error", receivedResponse.Err)
}

func TestGetUserFavoritesSuccess(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		user_favorites.UserFavorites{
			UserID:    1,
			MoviesIDs: []int{},
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "GET")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	GetUserFavorites(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse user_favorites.UserFavorites
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.NotNil(t, receivedResponse)
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, 1, len(receivedResponse.MoviesIDs))
	assert.EqualValues(t, 1, receivedResponse.MoviesIDs[0])
}

func TestGetUserFavoritesInvalidUserID(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		user_favorites.UserFavorites{
			UserID:    1,
			MoviesIDs: []int{},
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "GET")

	GetUserFavorites(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, "User ID should be a number", receivedResponse.Message)
	assert.EqualValues(t, "bad_request", receivedResponse.Err)
	assert.EqualValues(t, http.StatusBadRequest, receivedResponse.Status)
}

func TestGetUserFavoritesFailure(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		user_favorites.UserFavorites{
			UserID:    1,
			MoviesIDs: []int{},
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "GET")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	users_service.UsersService.(*usersServiceMock).canGetFavorites = false

	GetUserFavorites(c)

	users_service.UsersService.(*usersServiceMock).canGetFavorites = true

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse *rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.NotNil(t, receivedResponse)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, "Error when trying to get user favorites", receivedResponse.Message)
	assert.EqualValues(t, http.StatusInternalServerError, receivedResponse.Status)
	assert.EqualValues(t, "internal_server_error", receivedResponse.Err)
}

func TestAddUserFavoritesSuccessMovieCached(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		movies.MovieInfo{
			Movie: tmdb.Movie{
				ID:    1,
				Title: "Example Movie Title",
			},
			CreatedAt: "",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	AddUserFavorite(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	receivedResponse := string(responseData[:])

	assert.EqualValues(t, "", receivedResponse)
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, false, movies_service.MoviesService.(*moviesServiceMock).addedMovie)
}

func TestAddUserFavoritesSuccessMovieNotCached(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		movies.MovieInfo{
			Movie: tmdb.Movie{
				ID:    1,
				Title: "Example Movie Title",
			},
			CreatedAt: "",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	movies_service.MoviesService.(*moviesServiceMock).hasMovieCached = false

	AddUserFavorite(c)

	movies_service.MoviesService.(*moviesServiceMock).hasMovieCached = true

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	receivedResponse := string(responseData[:])

	assert.EqualValues(t, "", receivedResponse)
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, true, movies_service.MoviesService.(*moviesServiceMock).addedMovie)
}

func TestAddUserFavoritesInvalidUserID(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		movies.MovieInfo{
			Movie: tmdb.Movie{
				ID:    1,
				Title: "Example Movie Title",
			},
			CreatedAt: "",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	AddUserFavorite(c)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, "User ID should be a number", receivedResponse.Message)
	assert.EqualValues(t, "bad_request", receivedResponse.Err)
	assert.EqualValues(t, http.StatusBadRequest, receivedResponse.Status)
}

func TestAddUserFavoritesInvalidJSONBody(t *testing.T) {
	exampleJsonReq, err := json.Marshal(`{"invalid_key": "true"}`)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	AddUserFavorite(c)

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	assert.EqualValues(t, "Invalid JSON body", receivedResponse.Message)
	assert.EqualValues(t, "bad_request", receivedResponse.Err)
	assert.EqualValues(t, http.StatusBadRequest, receivedResponse.Status)
}

func TestAddUserFavoritesGetMovieError(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		movies.MovieInfo{
			Movie: tmdb.Movie{
				ID:    1,
				Title: "Example Movie Title",
			},
			CreatedAt: "",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	movies_service.MoviesService.(*moviesServiceMock).canGetMovie = false

	AddUserFavorite(c)

	movies_service.MoviesService.(*moviesServiceMock).canGetMovie = true

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, "Error when trying to get movie", receivedResponse.Message)
	assert.EqualValues(t, "internal_server_error", receivedResponse.Err)
	assert.EqualValues(t, http.StatusInternalServerError, receivedResponse.Status)
}

func TestAddUserFavoritesAddMovieErrorWhenNotCached(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		movies.MovieInfo{
			Movie: tmdb.Movie{
				ID:    1,
				Title: "Example Movie Title",
			},
			CreatedAt: "",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	movies_service.MoviesService.(*moviesServiceMock).hasMovieCached = false
	movies_service.MoviesService.(*moviesServiceMock).canAddMovie = false

	AddUserFavorite(c)

	movies_service.MoviesService.(*moviesServiceMock).hasMovieCached = true
	movies_service.MoviesService.(*moviesServiceMock).canAddMovie = true

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, "Error when trying to add movie", receivedResponse.Message)
	assert.EqualValues(t, "internal_server_error", receivedResponse.Err)
	assert.EqualValues(t, http.StatusInternalServerError, receivedResponse.Status)
}

func TestAddUserFavoritesAddUserFavoriteError(t *testing.T) {
	exampleJsonReq, err := json.Marshal(
		movies.MovieInfo{
			Movie: tmdb.Movie{
				ID:    1,
				Title: "Example Movie Title",
			},
			CreatedAt: "",
		},
	)
	if err != nil {
		panic(err)
	}

	w := PrepareTest(exampleJsonReq, "POST")

	c.Params = append(c.Params, gin.Param{Key: "user_id", Value: "1"})

	users_service.UsersService.(*usersServiceMock).canAddFavorite = false

	AddUserFavorite(c)

	users_service.UsersService.(*usersServiceMock).canAddFavorite = false

	c.Params = make([]gin.Param, 0)

	responseData, _ := ioutil.ReadAll(w.Body)

	var receivedResponse rest_errors.RestErr
	err = json.Unmarshal(responseData, &receivedResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, "Error when trying to add user favorite", receivedResponse.Message)
	assert.EqualValues(t, "internal_server_error", receivedResponse.Err)
	assert.EqualValues(t, http.StatusInternalServerError, receivedResponse.Status)
}
