package users

import (
	"context"
	"fmt"

	"github.com/ericbg27/top10movies-api/src/datasources/database"
	user_queries "github.com/ericbg27/top10movies-api/src/queries/users"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

func (user User) Get(db database.DatabaseClient) (UserInterface, *rest_errors.RestErr) {
	savedUser := user

	result, err := db.QueryRow(context.Background(), user_queries.QueryGetUser, user.Email)
	if err != nil {
		logger.Error("Error when trying to get user in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	err = result.Scan(&savedUser.ID, &savedUser.FirstName, &savedUser.Status, &savedUser.Password)
	if err != nil {
		logger.Error("Error when trying to get user in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	return savedUser, nil
}

func (user User) GetById(db database.DatabaseClient) (UserInterface, *rest_errors.RestErr) {
	savedUser := user

	result, err := db.QueryRow(context.Background(), user_queries.QueryGetUserById, user.ID)
	if err != nil {
		logger.Error("Error when trying to get user by id in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	err = result.Scan(&savedUser.FirstName, &savedUser.LastName, &savedUser.Email, &savedUser.Status, &savedUser.Password)
	if err != nil {
		logger.Error("Error when trying to get user by id in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	return savedUser, nil
}

func (user User) Save(db database.DatabaseClient) *rest_errors.RestErr {
	result, err := db.Exec(context.Background(), user_queries.QueryInsertUser, user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if err != nil {
		logger.Error("Error when trying to save user in database", err)
		return rest_errors.NewInternalServerError("Error when trying to save user")
	}

	logger.Info(fmt.Sprintf("Saved user in the database. Rows affected: %d", result.RowsAffected()))

	return nil
}

func (user User) Update(newUser UserInterface, isPartial bool, db database.DatabaseClient) (UserInterface, *rest_errors.RestErr) {
	toUpdateUser := newUser.(User)

	if isPartial {
		if toUpdateUser.FirstName != "" {
			user.FirstName = toUpdateUser.FirstName
		}
		if toUpdateUser.LastName != "" {
			user.LastName = toUpdateUser.LastName
		}
		if toUpdateUser.Email != "" {
			user.Email = toUpdateUser.Email
		}
	} else {
		user.FirstName = toUpdateUser.FirstName
		user.LastName = toUpdateUser.LastName
		user.Email = toUpdateUser.Email
	}

	validatedUser, validateErr := user.Validate()
	if validateErr != nil {
		return nil, validateErr
	}
	user = validatedUser.(User)

	result, err := db.Exec(context.Background(), user_queries.QueryUpdateUser, user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		logger.Error("Error when trying to update user in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to update user")
	}

	logger.Info(fmt.Sprintf("Updated user in the database. Rows affected: %d", result.RowsAffected()))

	return user, nil
}

func (user User) Delete(db database.DatabaseClient) *rest_errors.RestErr {
	result, err := db.Exec(context.Background(), user_queries.QueryDeleteUser, user.ID)
	if err != nil {
		logger.Error("Error when trying to delete user in database", err)
		return rest_errors.NewInternalServerError("Error when trying to delete user")
	}

	logger.Info(fmt.Sprintf("Deleted user in the database. Rows affected: %d", result.RowsAffected()))

	return nil
}

func (user User) Search(db database.DatabaseClient) ([]UserInterface, *rest_errors.RestErr) {
	result, err := db.Query(context.Background(), user_queries.QuerySearchUser, user.FirstName, user.LastName)
	if err != nil {
		logger.Error("Error when trying to search user in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to search user")
	}

	var foundUsers []UserInterface
	for result.Next() {
		var searchedUser User

		err = result.Scan(&searchedUser.ID, &searchedUser.FirstName, &searchedUser.LastName, &searchedUser.Email)
		if err != nil {
			logger.Error("Error when trying to search user in database", err)
			return nil, rest_errors.NewInternalServerError("Error when trying to search user")
		}

		foundUsers = append(foundUsers, searchedUser)
	}

	return foundUsers, nil
}
