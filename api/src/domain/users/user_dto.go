package users

import (
	"net/mail"
	"strings"

	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

const (
	StatusActive = "active"
)

type User struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

func (user *User) Validate() *rest_errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	if user.FirstName == "" || user.LastName == "" {
		return rest_errors.NewBadRequestError("First and last name fields cannot be empty")
	}

	user.Email = strings.TrimSpace(user.Email)
	if user.Email == "" {
		return rest_errors.NewBadRequestError("Invalid email address")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return rest_errors.NewBadRequestError("Invalid email address")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return rest_errors.NewBadRequestError("Invalid password")
	}

	return nil
}
