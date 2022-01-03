package users

import (
	"net/http"
	"testing"

	"github.com/ericbg27/top10movies-api/src/datasources/database"
	database_mock "github.com/ericbg27/top10movies-api/src/mocks/database"
	"github.com/stretchr/testify/assert"
)

var (
	db database.DatabaseClient
)

func TestGetSuccess(t *testing.T) {
	var user User

	result, err := user.Get(db)

	fetchedUser := result.(User)

	assert.Nil(t, err)
	assert.EqualValues(t, int64(1), fetchedUser.ID)
	assert.EqualValues(t, "John", fetchedUser.FirstName)
	assert.EqualValues(t, "Doe", fetchedUser.LastName)
	assert.EqualValues(t, "johndoe@gmail.com", fetchedUser.Email)
	assert.EqualValues(t, "1234", fetchedUser.Password)
}

func TestGetQueryRowError(t *testing.T) {
	var user User

	db.(*database_mock.DatabaseClientMock).CanQueryRow = false

	result, err := user.Get(db)

	db.(*database_mock.DatabaseClientMock).CanQueryRow = true

	assert.Nil(t, result)
	assert.EqualValues(t, "Error when trying to get user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestGetScanError(t *testing.T) {
	var user User

	db.(*database_mock.DatabaseClientMock).CanScanResults = false

	result, err := user.Get(db)

	db.(*database_mock.DatabaseClientMock).CanScanResults = true

	assert.Nil(t, result)
	assert.EqualValues(t, "Error when trying to get user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestGetByIdSuccess(t *testing.T) {
	var user User

	result, err := user.GetById(db)

	fetchedUser := result.(User)

	assert.Nil(t, err)
	assert.EqualValues(t, int64(1), fetchedUser.ID)
	assert.EqualValues(t, "John", fetchedUser.FirstName)
	assert.EqualValues(t, "Doe", fetchedUser.LastName)
	assert.EqualValues(t, "johndoe@gmail.com", fetchedUser.Email)
	assert.EqualValues(t, "1234", fetchedUser.Password)
}

func TestGetByIdQueryRowError(t *testing.T) {
	var user User

	db.(*database_mock.DatabaseClientMock).CanQueryRow = false

	result, err := user.Get(db)

	db.(*database_mock.DatabaseClientMock).CanQueryRow = true

	assert.Nil(t, result)
	assert.EqualValues(t, "Error when trying to get user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestGetByIdScanError(t *testing.T) {
	var user User

	db.(*database_mock.DatabaseClientMock).CanScanResults = false

	result, err := user.Get(db)

	db.(*database_mock.DatabaseClientMock).CanScanResults = true

	assert.Nil(t, result)
	assert.EqualValues(t, "Error when trying to get user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestSaveSuccess(t *testing.T) {
	var user User

	err := user.Save(db)

	assert.Nil(t, err)
}

func TestSaveExecError(t *testing.T) {
	var user User

	db.(*database_mock.DatabaseClientMock).CanExec = false

	err := user.Save(db)

	db.(*database_mock.DatabaseClientMock).CanExec = true

	assert.EqualValues(t, "Error when trying to save user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestUpdateSuccess(t *testing.T) {
	currentUser := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@mail.com",
		Password:  "1234",
	}

	updateUser := User{
		FirstName: "Johnn",
		LastName:  "Doee",
		Email:     "johnndoee@mail.com",
	}

	result, err := currentUser.Update(updateUser, false, db)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	updatedUser := result.(User)

	assert.EqualValues(t, "Johnn", updatedUser.FirstName)
	assert.EqualValues(t, "Doee", updatedUser.LastName)
	assert.EqualValues(t, "johnndoee@mail.com", updatedUser.Email)
	assert.EqualValues(t, updatedUser.ID, 1)
	assert.EqualValues(t, "1234", updatedUser.Password)
}

func TestUpdatePartialSuccess(t *testing.T) {
	currentUser := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@mail.com",
		Password:  "1234",
	}

	// First attempt
	updateUser := User{
		FirstName: "",
		LastName:  "Doee",
		Email:     "johnndoee@mail.com",
	}

	result, err := currentUser.Update(updateUser, true, db)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	updatedUser := result.(User)

	assert.EqualValues(t, "John", updatedUser.FirstName)
	assert.EqualValues(t, "Doee", updatedUser.LastName)
	assert.EqualValues(t, "johnndoee@mail.com", updatedUser.Email)
	assert.EqualValues(t, updatedUser.ID, 1)
	assert.EqualValues(t, "1234", updatedUser.Password)

	// Second attempt
	updateUser.FirstName = "Johnn"
	updateUser.LastName = ""

	result, err = currentUser.Update(updateUser, true, db)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	updatedUser = result.(User)

	assert.EqualValues(t, "Johnn", updatedUser.FirstName)
	assert.EqualValues(t, "Doe", updatedUser.LastName)
	assert.EqualValues(t, "johnndoee@mail.com", updatedUser.Email)
	assert.EqualValues(t, updatedUser.ID, 1)
	assert.EqualValues(t, "1234", updatedUser.Password)

	// Third attempt
	updateUser.FirstName = "Johnn"
	updateUser.LastName = "Doee"
	updateUser.Email = ""

	result, err = currentUser.Update(updateUser, true, db)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	updatedUser = result.(User)

	assert.EqualValues(t, "Johnn", updatedUser.FirstName)
	assert.EqualValues(t, "Doee", updatedUser.LastName)
	assert.EqualValues(t, "johndoe@mail.com", updatedUser.Email)
	assert.EqualValues(t, updatedUser.ID, 1)
	assert.EqualValues(t, "1234", updatedUser.Password)
}

func TestUpdateInvalidUserUpdate(t *testing.T) {
	currentUser := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@mail.com",
		Password:  "1234",
	}

	// First attempt
	updateUser := User{
		FirstName: "",
		LastName:  "Doee",
		Email:     "johnndoee@mail.com",
	}

	result, err := currentUser.Update(updateUser, false, db)

	assert.Nil(t, result)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "bad_request", err.Err)
	assert.EqualValues(t, "First and last name fields cannot be empty", err.Message)

	// Second attempt
	updateUser.FirstName = "Johnn"
	updateUser.LastName = ""

	result, err = currentUser.Update(updateUser, false, db)

	assert.Nil(t, result)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "bad_request", err.Err)
	assert.EqualValues(t, "First and last name fields cannot be empty", err.Message)

	// Third attempt
	updateUser.LastName = "Doee"
	updateUser.Email = ""

	result, err = currentUser.Update(updateUser, false, db)

	assert.Nil(t, result)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "bad_request", err.Err)
	assert.EqualValues(t, "Invalid email address", err.Message)

	// Fourth attempt
	updateUser.Email = "invalidemail.com"

	assert.Nil(t, result)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "bad_request", err.Err)
	assert.EqualValues(t, "Invalid email address", err.Message)
}

func TestUpdateExecError(t *testing.T) {
	currentUser := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@mail.com",
		Password:  "1234",
	}

	// First attempt
	updateUser := User{
		FirstName: "Johnn",
		LastName:  "Doee",
		Email:     "johnndoee@mail.com",
	}

	db.(*database_mock.DatabaseClientMock).CanExec = false

	result, err := currentUser.Update(updateUser, false, db)

	db.(*database_mock.DatabaseClientMock).CanExec = true

	assert.Nil(t, result)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
	assert.EqualValues(t, "Error when trying to update user", err.Message)
}

func TestDeleteSuccess(t *testing.T) {
	var user User

	err := user.Delete(db)

	assert.Nil(t, err)
}

func TestDeleteExecError(t *testing.T) {
	var user User

	db.(*database_mock.DatabaseClientMock).CanExec = false

	err := user.Delete(db)

	db.(*database_mock.DatabaseClientMock).CanExec = true

	assert.EqualValues(t, "Error when trying to delete user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}

func TestSearchSuccess(t *testing.T) {
	var user User

	results, err := user.Search(db)

	var usersFetched []User

	for _, result := range results {
		usersFetched = append(usersFetched, result.(User))
	}

	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(usersFetched))
	assert.EqualValues(t, int64(1), usersFetched[0].ID)
	assert.EqualValues(t, "John", usersFetched[0].FirstName)
	assert.EqualValues(t, "Doe", usersFetched[0].LastName)
	assert.EqualValues(t, "johndoe@gmail.com", usersFetched[0].Email)
	assert.EqualValues(t, "", usersFetched[0].Password)
	assert.EqualValues(t, int64(2), usersFetched[1].ID)
	assert.EqualValues(t, "Josh", usersFetched[1].FirstName)
	assert.EqualValues(t, "Davis", usersFetched[1].LastName)
	assert.EqualValues(t, "joshdavis@gmail.com", usersFetched[1].Email)
	assert.EqualValues(t, "", usersFetched[1].Password)
}

func TestSearchQueryError(t *testing.T) {
	var user User

	db.(*database_mock.DatabaseClientMock).CanQuery = false

	result, err := user.Search(db)

	db.(*database_mock.DatabaseClientMock).CanQuery = true

	assert.Nil(t, result)
	assert.EqualValues(t, "Error when trying to search user", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Err)
}
