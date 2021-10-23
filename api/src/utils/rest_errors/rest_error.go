package rest_errors

import (
	"fmt"
	"net/http"
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Err     string `json:"error"`
}

const (
	errorMessage              = "message: %s - status: %d - error: %s"
	badRequestString          = "bad_request"
	notFoundString            = "not_found"
	internalServerErrorString = "internal_server_error"
	unauthorizedString        = "unauthorized"
)

func (r RestErr) Error() string {
	return fmt.Sprintf(errorMessage, r.Message, r.Status, r.Err)
}

func NewRestError(message string, status int, err string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  status,
		Err:     err,
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Err:     badRequestString,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Err:     notFoundString,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Err:     internalServerErrorString,
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusUnauthorized,
		Err:     unauthorizedString,
	}
}
