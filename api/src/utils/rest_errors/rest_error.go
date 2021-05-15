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

func (r RestErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s", r.Message, r.Status, r.Err)
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
		Err:     "bad_request",
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Err:     "not_found",
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Err:     "internal_server_error",
	}
}
