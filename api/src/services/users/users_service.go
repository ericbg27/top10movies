package users_service

import (
	"github.com/ericbg27/top10movies-api/src/datasources/database"
	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
	"github.com/ericbg27/top10movies-api/src/domain/user_favorites"
	"github.com/ericbg27/top10movies-api/src/domain/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

type usersService struct {
	db          database.DatabaseClient
	redisClient redisdb.RedisClient
}

type UsersServiceInterface interface {
	CreateUser(users.UserInterface) (users.UserInterface, *rest_errors.RestErr)
	GetUser(users.UserInterface) (users.UserInterface, *rest_errors.RestErr)
	UpdateUser(users.UserInterface, bool) (users.UserInterface, *rest_errors.RestErr)
	DeleteUser(users.UserInterface) *rest_errors.RestErr
	GetUserFavorites(user_favorites.UserFavoritesInterface) (user_favorites.UserFavoritesInterface, map[int]bool, *rest_errors.RestErr)
	AddUserFavorite(user_favorites.UserFavoritesInterface) *rest_errors.RestErr
	SearchUser(users.UserInterface) ([]users.UserInterface, *rest_errors.RestErr)
}

const (
	QueryParam = "query"
)

func NewUsersService(db database.DatabaseClient, rc redisdb.RedisClient) *usersService {
	m := &usersService{
		db:          db,
		redisClient: rc,
	}

	return m
}

func (s *usersService) GetUser(user users.UserInterface) (users.UserInterface, *rest_errors.RestErr) {
	var savedUser users.UserInterface
	var err *rest_errors.RestErr

	if savedUser, err = user.Get(s.db); err != nil {
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

	if err = validatedUser.Save(s.db); err != nil {
		return nil, err
	}

	return validatedUser, nil
}

func (s *usersService) UpdateUser(user users.UserInterface, isPartial bool) (users.UserInterface, *rest_errors.RestErr) {
	var currentUser users.UserInterface
	var err *rest_errors.RestErr

	if currentUser, err = user.GetById(s.db); err != nil {
		return nil, err
	}

	var updatedUser users.UserInterface
	if updatedUser, err = currentUser.Update(user, isPartial, s.db); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *usersService) DeleteUser(user users.UserInterface) *rest_errors.RestErr {
	var currentUser users.UserInterface
	var err *rest_errors.RestErr

	if currentUser, err = user.GetById(s.db); err != nil {
		return err
	}

	if err = currentUser.Delete(s.db); err != nil {
		return err
	}

	return nil
}

func (s *usersService) GetUserFavorites(userFavorites user_favorites.UserFavoritesInterface) (user_favorites.UserFavoritesInterface, map[int]bool, *rest_errors.RestErr) {
	var currentUserFavorites user_favorites.UserFavoritesInterface
	var cachedIds map[int]bool
	var err *rest_errors.RestErr

	if currentUserFavorites, cachedIds, err = userFavorites.GetFavorites(s.db, s.redisClient); err != nil {
		return nil, nil, err
	}

	return currentUserFavorites, cachedIds, nil
}

func (s *usersService) AddUserFavorite(userFavorites user_favorites.UserFavoritesInterface) *rest_errors.RestErr {
	if err := userFavorites.AddFavorite(s.db); err != nil {
		return err
	}

	return nil
}

func (s *usersService) SearchUser(userToSearch users.UserInterface) ([]users.UserInterface, *rest_errors.RestErr) {
	usersFound, searchErr := userToSearch.Search(s.db)
	if searchErr != nil {
		return nil, searchErr
	}

	return usersFound, nil
}
