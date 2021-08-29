package users

import (
	"fmt"

	"github.com/ericbg27/top10movies-api/src/datasources/postgresql/users_db"
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
)

func (user User) Get() (UserInterface, *rest_errors.RestErr) {
	savedUser := user
	_, err := users_db.Client.Prepare(queryGetUserName, queryGetUser)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	result := users_db.Client.QueryRow(queryGetUserName, user.Email)
	err = result.Scan(&savedUser.ID, &savedUser.FirstName, &savedUser.Status, &savedUser.Password)
	if err != nil {
		logger.Error("Error when trying to get user in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	return savedUser, nil
}

func (user User) GetById() (UserInterface, *rest_errors.RestErr) {
	savedUser := user
	_, err := users_db.Client.Prepare(queryGetUserByIdName, queryGetUserById)
	if err != nil {
		logger.Error("Error when trying to prepare get user by id statement", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	result := users_db.Client.QueryRow(queryGetUserByIdName, user.ID)
	err = result.Scan(&savedUser.FirstName, &savedUser.LastName, &savedUser.Email, &savedUser.Status, &savedUser.Password)
	if err != nil {
		logger.Error("Error when trying to get user by id in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get user")
	}

	return savedUser, nil
}

func (user User) Save() *rest_errors.RestErr {
	_, err := users_db.Client.Prepare(queryInsertUserName, queryInsertUser)
	if err != nil {
		logger.Error("Error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("Error when trying to save user")
	}

	result, err := users_db.Client.Exec(queryInsertUserName, user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
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

	_, err := users_db.Client.Prepare(queryUpdateUserName, queryUpdateUser)
	if err != nil {
		logger.Error("Error when trying to prepare update user statement", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to update user")
	}

	result, err := users_db.Client.Exec(queryUpdateUserName, user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		logger.Error("Error when trying to update user in database", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to update user")
	}

	logger.Info(fmt.Sprintf("Updated user in the database. Rows affected: %d", result.RowsAffected()))

	return user, nil
}
