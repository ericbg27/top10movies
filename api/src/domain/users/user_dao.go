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
)

func (user *User) Save() *rest_errors.RestErr {
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
