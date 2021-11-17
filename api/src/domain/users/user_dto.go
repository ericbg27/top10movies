package users

import (
	"net/mail"
	"strings"

	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

const (
	StatusActive = "active"
)

type UserInterface interface {
	Validate() (UserInterface, *rest_errors.RestErr)
	Get() (UserInterface, *rest_errors.RestErr)
	GetById() (UserInterface, *rest_errors.RestErr)
	Save() *rest_errors.RestErr
	Update(newUser UserInterface, isPartial bool) (UserInterface, *rest_errors.RestErr)
	Delete() *rest_errors.RestErr
	Search() ([]UserInterface, *rest_errors.RestErr)
}

type User struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

func (user User) Validate() (UserInterface, *rest_errors.RestErr) {
	validatedUser := user
	validatedUser.FirstName = strings.TrimSpace(validatedUser.FirstName)
	validatedUser.LastName = strings.TrimSpace(validatedUser.LastName)
	if validatedUser.FirstName == "" || validatedUser.LastName == "" {
		return nil, rest_errors.NewBadRequestError("First and last name fields cannot be empty")
	}

	validatedUser.Email = strings.TrimSpace(validatedUser.Email)
	if validatedUser.Email == "" {
		return nil, rest_errors.NewBadRequestError("Invalid email address")
	}

	_, err := mail.ParseAddress(validatedUser.Email)
	if err != nil {
		return nil, rest_errors.NewBadRequestError("Invalid email address")
	}

	validatedUser.Password = strings.TrimSpace(validatedUser.Password)
	if validatedUser.Password == "" {
		return nil, rest_errors.NewBadRequestError("Invalid password")
	}

	return validatedUser, nil
}
