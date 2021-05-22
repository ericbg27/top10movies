package users_service

import (
	"net/http"
	"os"
	"testing"

	"github.com/ericbg27/top10movies-api/src/domain/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/stretchr/testify/assert"
)

type userMock struct {
	valid   bool
	canGet  bool
	canSave bool
	name    string
}

func (u userMock) Validate() (users.UserInterface, *rest_errors.RestErr) {
	validatedUser := u

	if validatedUser.valid == false {
		return nil, rest_errors.NewBadRequestError("Invalid user")
	}

	return validatedUser, nil
}

func (u userMock) Get() (users.UserInterface, *rest_errors.RestErr) {
	savedUser := u

	if savedUser.canGet == false {
		return nil, rest_errors.NewInternalServerError("Failed to get user")
	}

	savedUser.name = "Test User"

	return savedUser, nil
}

func (u userMock) Save() *rest_errors.RestErr {
	if u.canSave == false {
		return rest_errors.NewInternalServerError("Failed to save user")
	}

	return nil
}

func TestMain(m *testing.M) {
	UsersService = &usersService{}
	os.Exit(m.Run())
}

func TestGetUserSuccess(t *testing.T) {
	var user userMock
	user.canGet = true

	result, err := UsersService.GetUser(user)

	savedUser := result.(userMock)

	assert.Nil(t, err)
	assert.NotNil(t, savedUser)
	assert.EqualValues(t, "Test User", savedUser.name)
	assert.EqualValues(t, "", user.name)
}

func TestGetUserFail(t *testing.T) {
	var user userMock
	user.canGet = false

	result, err := UsersService.GetUser(user)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Failed to get user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestCreateUserSuccess(t *testing.T) {
	var user userMock
	user.valid = true
	user.canSave = true
	user.name = "User to create"

	result, err := UsersService.CreateUser(user)

	createdUser := result.(userMock)

	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	assert.EqualValues(t, "User to create", createdUser.name)
}

func TestCreateUserInvalidUser(t *testing.T) {
	var user userMock
	user.valid = false

	result, err := UsersService.CreateUser(user)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Invalid user", err.Message)
	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "bad_request", err.Err)
}

func TestCreateUserFail(t *testing.T) {
	var user userMock
	user.valid = true
	user.canSave = false

	result, err := UsersService.CreateUser(user)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Failed to save user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}
