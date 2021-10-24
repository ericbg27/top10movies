package rest_errors

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	forbiddenString     = "forbidden"
	statusCreatedString = "status_created"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestError(t *testing.T) {
	var err RestErr = RestErr{
		Message: "Default Error",
		Status:  http.StatusForbidden,
		Err:     forbiddenString,
	}

	result := err.Error()
	expected := fmt.Sprintf(errorMessage, "Default Error", http.StatusForbidden, forbiddenString)

	assert.Equal(t, expected, result)
}

func TestNewRestError(t *testing.T) {
	restErr := NewRestError("Test Message", http.StatusCreated, statusCreatedString)

	assert.EqualValues(t, "Test Message", restErr.Message)
	assert.EqualValues(t, http.StatusCreated, restErr.Status)
	assert.EqualValues(t, statusCreatedString, restErr.Err)
}

func TestNewBadRequestError(t *testing.T) {
	badRequestErr := NewBadRequestError("Bad Request")

	assert.EqualValues(t, "Bad Request", badRequestErr.Message)
	assert.EqualValues(t, http.StatusBadRequest, badRequestErr.Status)
	assert.EqualValues(t, badRequestString, badRequestErr.Err)
}

func TestNewNotFoundError(t *testing.T) {
	notFoundErr := NewNotFoundError("Not Found")

	assert.EqualValues(t, "Not Found", notFoundErr.Message)
	assert.EqualValues(t, http.StatusNotFound, notFoundErr.Status)
	assert.EqualValues(t, notFoundString, notFoundErr.Err)
}

func TestNewInternalServerError(t *testing.T) {
	internalServerErr := NewInternalServerError("Internal Server Error")

	assert.EqualValues(t, "Internal Server Error", internalServerErr.Message)
	assert.EqualValues(t, http.StatusInternalServerError, internalServerErr.Status)
	assert.EqualValues(t, internalServerErrorString, internalServerErr.Err)
}

func TestNewUnauthorizedError(t *testing.T) {
	unauthorizedErr := NewUnauthorizedError("Unauthorized")

	assert.EqualValues(t, "Unauthorized", unauthorizedErr.Message)
	assert.EqualValues(t, http.StatusUnauthorized, unauthorizedErr.Status)
	assert.EqualValues(t, unauthorizedString, unauthorizedErr.Err)
}
