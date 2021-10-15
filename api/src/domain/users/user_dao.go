package users

import (
	"context"
	"fmt"

	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/db"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

const (
	queryInsertUser     = "INSERT INTO users (first_name,last_name,email,date_created,status,password) VALUES ($1,$2,$3,$4,$5,$6);"
	queryInsertUserName = "insert-user-query"

	queryGetUser     = "SELECT id, first_name, status, password FROM users WHERE email=$1;"
	queryGetUserName = "get-user-query"

	queryGetUserById     = "SELECT first_name, last_name, email, status, password FROM users WHERE id=$1;"
	queryGetUserByIdName = "get-user-by-id-query"

	queryUpdateUser     = "UPDATE users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
	queryUpdateUserName = "update-user-query"

	queryDeleteUser     = "DELETE FROM users WHERE id=$1;"
	queryDeleteUserName = "delete-user-query"
)

func (user User) Get() (UserInterface, *rest_errors.RestErr) {
	savedUser := user

	result := db.Client.QueryRow(context.Background(), queryGetUser, user.Email)
	err := result.Scan(&savedUser.ID, &savedUser.FirstName, &savedUser.Status, &savedUser.Password)
	if err != nil {
		logger.Error("Error when trying to get user in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	return savedUser, nil
}

func (user User) GetById() (UserInterface, *rest_errors.RestErr) {
	savedUser := user

	result := db.Client.QueryRow(context.Background(), queryGetUserById, user.ID)
	err := result.Scan(&savedUser.FirstName, &savedUser.LastName, &savedUser.Email, &savedUser.Status, &savedUser.Password)
	if err != nil {
		logger.Error("Error when trying to get user by id in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	return savedUser, nil
}

func (user User) Save() *rest_errors.RestErr {
	result, err := db.Client.Exec(context.Background(), queryInsertUser, user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if err != nil {
		logger.Error("Error when trying to save user in database", err)
		return rest_errors.NewInternalServerError("Error when trying to save user")
	}

	logger.Info(fmt.Sprintf("Saved user in the database. Rows affected: %d", result.RowsAffected()))

	return nil
}

func (user User) Update(newUser UserInterface, isPartial bool) (UserInterface, *rest_errors.RestErr) {
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

	result, err := db.Client.Exec(context.Background(), queryUpdateUser, user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		logger.Error("Error when trying to update user in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to update user")
	}

	logger.Info(fmt.Sprintf("Updated user in the database. Rows affected: %d", result.RowsAffected()))

	return user, nil
}

func (user User) Delete() *rest_errors.RestErr {
	result, err := db.Client.Exec(context.Background(), queryDeleteUser, user.ID)
	if err != nil {
		logger.Error("Error when trying to delete user in database", err)
		return rest_errors.NewInternalServerError("Error when trying to delete user")
	}

	logger.Info(fmt.Sprintf("Deleted user in the database. Rows affected: %d", result.RowsAffected()))

	return nil
}
