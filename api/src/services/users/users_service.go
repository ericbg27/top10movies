package users_service

import (
	"github.com/ericbg27/top10movies-api/src/domain/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

type usersService struct{}

type usersServiceInterface interface {
	CreateUser(users.UserInterface) (users.UserInterface, *rest_errors.RestErr)
	GetUser(users.UserInterface) (users.UserInterface, *rest_errors.RestErr)
}

var (
	UsersService usersServiceInterface = &usersService{}
)

func (s *usersService) GetUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	var savedUser users.UserInterface
	var err *rest_errors.RestErr

	if savedUser, err = user.Get(); err != nil {
		return nil, err
	}

	return savedUser, nil
}

func (s *usersService) CreateUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	var validatedUser users.UserInterface
	var err *rest_errors.RestErr

	if validatedUser, err = user.Validate(); err != nil {
		return nil, err
	}

	if err = validatedUser.Save(); err != nil {
		return nil, err
	}

	return validatedUser, nil
}
