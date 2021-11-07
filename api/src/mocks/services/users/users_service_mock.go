package users_service

import (
	"github.com/ericbg27/top10movies-api/src/domain/user_favorites"
	"github.com/ericbg27/top10movies-api/src/domain/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/ryanbradynd05/go-tmdb"
)

var (
	MockDb   map[string]string
	MockDbID map[int64]users.User
	Now      string
)

type UsersServiceMock struct {
	CanGetFavorites bool
	CanAddFavorite  bool
	FavoriteCached  bool
}

func (u *UsersServiceMock) CreateUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	usr := user.(users.User)
	if _, ok := MockDb[usr.Email]; ok {
		return nil, rest_errors.NewInternalServerError("Error when trying to save user")
	}

	usr.DateCreated = Now
	usr.ID = 2

	return usr, nil
}

func (u *UsersServiceMock) GetUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	usr := user.(users.User)
	if savedPassword, ok := MockDb[usr.Email]; ok {
		savedUser := users.User{
			ID:       usr.ID,
			Email:    usr.Email,
			Password: savedPassword,
		}

		return savedUser, nil
	}

	return nil, rest_errors.NewNotFoundError("User not found")
}

func (u *UsersServiceMock) UpdateUser(user users.UserInterface, isPartial bool) (users.UserInterface, *rest_errors.RestErr) {
	newUser := user.(users.User)

	currentUser, ok := MockDbID[newUser.ID]
	if !ok {
		return nil, rest_errors.NewInternalServerError("Error when trying to update user")
	}

	if isPartial {
		if newUser.FirstName == "" {
			newUser.FirstName = currentUser.FirstName
		}
		if newUser.LastName == "" {
			newUser.LastName = currentUser.LastName
		}
		if newUser.Email == "" {
			newUser.Email = currentUser.Email
		}
	}

	return newUser, nil
}

func (u *UsersServiceMock) DeleteUser(user users.UserInterface) *rest_errors.RestErr {
	usr := user.(users.User)

	_, ok := MockDbID[usr.ID]
	if !ok {
		return rest_errors.NewInternalServerError("Error when trying to delete user")
	}

	return nil
}

func (u *UsersServiceMock) GetUserFavorites(userFavs user_favorites.UserFavoritesInterface) (user_favorites.UserFavoritesInterface, map[int]bool, *rest_errors.RestErr) {
	userFavorites := userFavs.(user_favorites.UserFavorites)

	if !u.CanGetFavorites {
		return nil, nil, rest_errors.NewInternalServerError("Error when trying to get user favorites")
	}

	userFavorites.MoviesIDs = append(userFavorites.MoviesIDs, 1)

	cacheMap := make(map[int]bool)
	if u.FavoriteCached {
		userFavorites.MoviesData = append(userFavorites.MoviesData, tmdb.Movie{
			ID: 1,
		})
		cacheMap[1] = true
	}

	return userFavorites, cacheMap, nil
}

func (u *UsersServiceMock) AddUserFavorite(userFavs user_favorites.UserFavoritesInterface) *rest_errors.RestErr {
	if !u.CanAddFavorite {
		return rest_errors.NewInternalServerError("Error when trying to add user favorite")
	}

	return nil
}

func (u *UsersServiceMock) SearchUser(userToSearch users.UserInterface) ([]users.UserInterface, *rest_errors.RestErr) {
	// TODO
	return nil, nil
}
