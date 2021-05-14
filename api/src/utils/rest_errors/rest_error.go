package rest_errors

import (
	"fmt"
	"net/http"
)

type restErr struct {
	message string `json:"message"`
	status  int    `json:"status"`
	err     string `json:"error"`
}

type RestErr interface {
	Message() string
	Status() int
	Error() string
}

func (r restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s", r.message, r.status, r.err)
}

func (r restErr) Message() string {
	return r.message
}

func (r restErr) Status() int {
	return r.status
}

func NewRestError(message string, status int, err string) RestErr {
	return restErr{
		message: message,
		status:  status,
		err:     err,
	}
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		message: message,
		status:  http.StatusBadRequest,
		err:     "bad_request",
	}
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		message: message,
		status:  http.StatusNotFound,
		err:     "not_found",
	}
}

func NewInternalServerError(message string) RestErr {
	return restErr{
		message: message,
		status:  http.StatusInternalServerError,
		err:     "internal_server_error",
	}
}
