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

	"github.com/ericbg27/top10movies-api/src/domain/users"
	users_service "github.com/ericbg27/top10movies-api/src/services/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	c      *gin.Context
	mockDb map[string]string
	now    string
)

type usersServiceMock struct{}

func (u *usersServiceMock) CreateUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	usr := user.(*users.User)
	if _, ok := mockDb[usr.Email]; ok {
		return nil, rest_errors.NewInternalServerError("Error when trying to save user")
	}

	usr.DateCreated = now
	usr.ID = 2

	return usr, nil
}

func (u *usersServiceMock) GetUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	usr := user.(*users.User)
	if savedPassword, ok := mockDb[usr.Email]; ok {
		savedUser := users.User{
			Email:    usr.Email,
			Password: savedPassword,
		}

		return &savedUser, nil
	}

	return nil, rest_errors.NewNotFoundError("User not found")
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

	gin.SetMode(gin.TestMode)

	now = time.Now().Format(layoutISO)

	users_service.UsersService = &usersServiceMock{}

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
