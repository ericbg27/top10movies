package users_service

import (
	"github.com/ericbg27/top10movies-api/src/domain/user_favorites"
	"github.com/ericbg27/top10movies-api/src/domain/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

type usersService struct{}

type usersServiceInterface interface {
	CreateUser(users.UserInterface) (users.UserInterface, *rest_errors.RestErr)
	GetUser(users.UserInterface) (users.UserInterface, *rest_errors.RestErr)
	UpdateUser(users.UserInterface, bool) (users.UserInterface, *rest_errors.RestErr)
	DeleteUser(users.UserInterface) *rest_errors.RestErr
	GetUserFavorites(user_favorites.UserFavoritesInterface) (user_favorites.UserFavoritesInterface, map[int]bool, *rest_errors.RestErr)
	AddUserFavorite(user_favorites.UserFavoritesInterface) *rest_errors.RestErr
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

func (s *usersService) UpdateUser(user users.UserInterface, isPartial bool) (users.UserInterface, *rest_errors.RestErr) {
	var currentUser users.UserInterface
	var err *rest_errors.RestErr

	if currentUser, err = user.GetById(); err != nil {
		return nil, err
	}

	var updatedUser users.UserInterface
	if updatedUser, err = currentUser.Update(user, isPartial); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *usersService) DeleteUser(user users.UserInterface) *rest_errors.RestErr {
	var currentUser users.UserInterface
	var err *rest_errors.RestErr

	if currentUser, err = user.GetById(); err != nil {
		return err
	}

	if err = currentUser.Delete(); err != nil {
		return err
	}

	return nil
}

func (s *usersService) GetUserFavorites(userFavorites user_favorites.UserFavoritesInterface) (user_favorites.UserFavoritesInterface, map[int]bool, *rest_errors.RestErr) {
	var currentUserFavorites user_favorites.UserFavoritesInterface
	var cachedIds map[int]bool
	var err *rest_errors.RestErr

	if currentUserFavorites, cachedIds, err = userFavorites.GetFavorites(); err != nil {
		return nil, nil, err
	}

	return currentUserFavorites, cachedIds, nil
}

func (s *usersService) AddUserFavorite(userFavorites user_favorites.UserFavoritesInterface) *rest_errors.RestErr {
	if err := userFavorites.AddFavorite(); err != nil {
		return err
	}

	return nil
}
