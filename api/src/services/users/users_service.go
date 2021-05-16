package users_service

import (
	"github.com/ericbg27/top10movies-api/src/domain/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

type usersService struct{}

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *rest_errors.RestErr)
	GetUser(users.User) (*users.User, *rest_errors.RestErr)
}

var (
	UsersService usersServiceInterface = &usersService{}
)

func (s *usersService) GetUser(user users.User) (*users.User, *rest_errors.RestErr) {
	if err := user.Get(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, *rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}
