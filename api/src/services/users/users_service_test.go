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
	valid     bool
	canGet    bool
	canSave   bool
	canUpdate bool
	canDelete bool
	FirstName string
	LastName  string
	Email     string
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

	savedUser.FirstName = "Test User"

	return savedUser, nil
}

func (u userMock) GetById() (users.UserInterface, *rest_errors.RestErr) {
	savedUser := u

	if savedUser.canGet == false {
		return nil, rest_errors.NewInternalServerError("Failed to get user by ID")
	}

	savedUser.FirstName = "Current Name"
	savedUser.LastName = "Current Last Name"
	savedUser.Email = "Current email"

	return savedUser, nil
}

func (u userMock) Save() *rest_errors.RestErr {
	if u.canSave == false {
		return rest_errors.NewInternalServerError("Failed to save user")
	}

	return nil
}

func (u userMock) Update(newUser users.UserInterface, isPartial bool) (users.UserInterface, *rest_errors.RestErr) {
	var validatedNewUser users.UserInterface
	var err *rest_errors.RestErr

	if validatedNewUser, err = newUser.Validate(); err != nil {
		return nil, err
	}

	if u.canUpdate == false {
		return nil, rest_errors.NewInternalServerError("Failed to update user")
	}

	toUpdateUser := validatedNewUser.(userMock)

	if isPartial {
		if toUpdateUser.FirstName != "" {
			u.FirstName = toUpdateUser.FirstName
		}
		if toUpdateUser.LastName != "" {
			u.LastName = toUpdateUser.LastName
		}
		if toUpdateUser.Email != "" {
			u.Email = toUpdateUser.Email
		}
	} else {
		u.FirstName = toUpdateUser.FirstName
		u.LastName = toUpdateUser.LastName
		u.Email = toUpdateUser.Email
	}

	return u, nil
}

func (u userMock) Delete() *rest_errors.RestErr {
	if u.canDelete == false {
		return rest_errors.NewInternalServerError("Failed to delete user")
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
	assert.EqualValues(t, "Test User", savedUser.FirstName)
	assert.EqualValues(t, "", user.FirstName)
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
	user.FirstName = "User to create"

	result, err := UsersService.CreateUser(user)

	createdUser := result.(userMock)

	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	assert.EqualValues(t, "User to create", createdUser.FirstName)
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

func TestUpdateUserSuccess(t *testing.T) {
	var user userMock
	user.valid = true
	user.canGet = true
	user.canUpdate = true
	user.FirstName = "Test Name"
	user.LastName = "Test Last Name"
	user.Email = "test@email.com"

	result, err := UsersService.UpdateUser(user, false)

	updatedUser := result.(userMock)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, "Test Name", updatedUser.FirstName)
	assert.EqualValues(t, "Test Last Name", updatedUser.LastName)
	assert.EqualValues(t, "test@email.com", updatedUser.Email)
}

func TestPartialUpdateUserSuccess(t *testing.T) {
	var user userMock
	user.valid = true
	user.canGet = true
	user.canUpdate = true
	user.FirstName = ""
	user.LastName = "Test Last Name"
	user.Email = "test@email.com"

	result, err := UsersService.UpdateUser(user, true)

	updatedUser := result.(userMock)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, "Current Name", updatedUser.FirstName)
	assert.EqualValues(t, "Test Last Name", updatedUser.LastName)
	assert.EqualValues(t, "test@email.com", updatedUser.Email)
}

func TestUpdateUserInvalidUserID(t *testing.T) {
	var user userMock
	user.valid = true
	user.canGet = false
	user.canUpdate = true

	result, err := UsersService.UpdateUser(user, false)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "Failed to get user by ID", err.Message)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestUpdateUserUpdateError(t *testing.T) {
	var user userMock
	user.valid = true
	user.canGet = true
	user.canUpdate = false

	result, err := UsersService.UpdateUser(user, false)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "Failed to update user", err.Message)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestUpdateUserInvalidUser(t *testing.T) {
	var user userMock
	user.valid = false
	user.canGet = true

	result, err := UsersService.UpdateUser(user, false)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Invalid user", err.Message)
	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "bad_request", err.Err)
}

func TestDeleteUserSuccess(t *testing.T) {
	var user userMock
	user.canGet = true
	user.canDelete = true

	err := UsersService.DeleteUser(user)

	assert.Nil(t, err)
}

func TestDeleteUserGetError(t *testing.T) {
	var user userMock
	user.canGet = false
	user.canDelete = true

	err := UsersService.DeleteUser(user)

	assert.NotNil(t, err)
	assert.EqualValues(t, "Failed to get user by ID", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestDeleteUserDeleteError(t *testing.T) {
	var user userMock
	user.canGet = true
	user.canDelete = false

	err := UsersService.DeleteUser(user)

	assert.NotNil(t, err)
	assert.EqualValues(t, "Failed to delete user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}
