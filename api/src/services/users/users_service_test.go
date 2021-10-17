package users_service

import (
	"net/http"
	"os"
	"testing"

	users_mock "github.com/ericbg27/top10movies-api/src/mocks/domain/users"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	UsersService = &usersService{}
	os.Exit(m.Run())
}

func TestGetUserSuccess(t *testing.T) {
	var user users_mock.UserMock
	user.CanGet = true

	result, err := UsersService.GetUser(user)

	savedUser := result.(users_mock.UserMock)

	assert.Nil(t, err)
	assert.NotNil(t, savedUser)
	assert.EqualValues(t, "Test User", savedUser.FirstName)
	assert.EqualValues(t, "", user.FirstName)
}

func TestGetUserFail(t *testing.T) {
	var user users_mock.UserMock
	user.CanGet = false

	result, err := UsersService.GetUser(user)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Failed to get user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestCreateUserSuccess(t *testing.T) {
	var user users_mock.UserMock
	user.Valid = true
	user.CanSave = true
	user.FirstName = "User to create"

	result, err := UsersService.CreateUser(user)

	createdUser := result.(users_mock.UserMock)

	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	assert.EqualValues(t, "User to create", createdUser.FirstName)
}

func TestCreateUserInvalidUser(t *testing.T) {
	var user users_mock.UserMock
	user.Valid = false

	result, err := UsersService.CreateUser(user)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Invalid user", err.Message)
	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "bad_request", err.Err)
}

func TestCreateUserFail(t *testing.T) {
	var user users_mock.UserMock
	user.Valid = true
	user.CanSave = false

	result, err := UsersService.CreateUser(user)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Failed to save user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestUpdateUserSuccess(t *testing.T) {
	var user users_mock.UserMock
	user.Valid = true
	user.CanGet = true
	user.CanUpdate = true
	user.FirstName = "Test Name"
	user.LastName = "Test Last Name"
	user.Email = "test@email.com"

	result, err := UsersService.UpdateUser(user, false)

	updatedUser := result.(users_mock.UserMock)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, "Test Name", updatedUser.FirstName)
	assert.EqualValues(t, "Test Last Name", updatedUser.LastName)
	assert.EqualValues(t, "test@email.com", updatedUser.Email)
}

func TestPartialUpdateUserSuccess(t *testing.T) {
	var user users_mock.UserMock
	user.Valid = true
	user.CanGet = true
	user.CanUpdate = true
	user.FirstName = ""
	user.LastName = "Test Last Name"
	user.Email = "test@email.com"

	result, err := UsersService.UpdateUser(user, true)

	updatedUser := result.(users_mock.UserMock)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, "Current Name", updatedUser.FirstName)
	assert.EqualValues(t, "Test Last Name", updatedUser.LastName)
	assert.EqualValues(t, "test@email.com", updatedUser.Email)
}

func TestUpdateUserInvalidUserID(t *testing.T) {
	var user users_mock.UserMock
	user.Valid = true
	user.CanGet = false
	user.CanUpdate = true

	result, err := UsersService.UpdateUser(user, false)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "Failed to get user by ID", err.Message)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestUpdateUserUpdateError(t *testing.T) {
	var user users_mock.UserMock
	user.Valid = true
	user.CanGet = true
	user.CanUpdate = false

	result, err := UsersService.UpdateUser(user, false)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "Failed to update user", err.Message)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestUpdateUserInvalidUser(t *testing.T) {
	var user users_mock.UserMock
	user.Valid = false
	user.CanGet = true

	result, err := UsersService.UpdateUser(user, false)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Invalid user", err.Message)
	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "bad_request", err.Err)
}

func TestDeleteUserSuccess(t *testing.T) {
	var user users_mock.UserMock
	user.CanGet = true
	user.CanDelete = true

	err := UsersService.DeleteUser(user)

	assert.Nil(t, err)
}

func TestDeleteUserGetError(t *testing.T) {
	var user users_mock.UserMock
	user.CanGet = false
	user.CanDelete = true

	err := UsersService.DeleteUser(user)

	assert.NotNil(t, err)
	assert.EqualValues(t, "Failed to get user by ID", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestDeleteUserDeleteError(t *testing.T) {
	var user users_mock.UserMock
	user.CanGet = true
	user.CanDelete = false

	err := UsersService.DeleteUser(user)

	assert.NotNil(t, err)
	assert.EqualValues(t, "Failed to delete user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}
