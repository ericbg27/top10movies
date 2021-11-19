package users

import (
	"github.com/ericbg27/top10movies-api/src/datasources/database"
	"github.com/ericbg27/top10movies-api/src/domain/users"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

type UserMock struct {
	Valid     bool
	CanGet    bool
	CanSave   bool
	CanUpdate bool
	CanDelete bool
	FirstName string
	LastName  string
	Email     string
}

// TODO: Is it a good idea for UserMock functions to user database mock structures? The input is already required

func (u UserMock) Validate() (users.UserInterface, *rest_errors.RestErr) {
	validatedUser := u

	if !validatedUser.Valid {
		return nil, rest_errors.NewBadRequestError("Invalid user")
	}

	return validatedUser, nil
}

func (u UserMock) Get(db database.DatabaseClient) (users.UserInterface, *rest_errors.RestErr) {
	savedUser := u

	if !savedUser.CanGet {
		return nil, rest_errors.NewInternalServerError("Failed to get user")
	}

	savedUser.FirstName = "Test User"

	return savedUser, nil
}

func (u UserMock) GetById(db database.DatabaseClient) (users.UserInterface, *rest_errors.RestErr) {
	savedUser := u

	if !savedUser.CanGet {
		return nil, rest_errors.NewInternalServerError("Failed to get user by ID")
	}

	savedUser.FirstName = "Current Name"
	savedUser.LastName = "Current Last Name"
	savedUser.Email = "Current email"

	return savedUser, nil
}

func (u UserMock) Save(db database.DatabaseClient) *rest_errors.RestErr {
	if !u.CanSave {
		return rest_errors.NewInternalServerError("Failed to save user")
	}

	return nil
}

func (u UserMock) Update(newUser users.UserInterface, isPartial bool, db database.DatabaseClient) (users.UserInterface, *rest_errors.RestErr) {
	var validatedNewUser users.UserInterface
	var err *rest_errors.RestErr

	if validatedNewUser, err = newUser.Validate(); err != nil {
		return nil, err
	}

	if !u.CanUpdate {
		return nil, rest_errors.NewInternalServerError("Failed to update user")
	}

	toUpdateUser := validatedNewUser.(UserMock)

	if isPartial {
		if toUpdateUser.FirstName != "" {
			u.FirstName = toUpdateUser.FirstName
		}
		if toUpdateUser.LastName != "" {
			u.LastName = toUpdateUser.LastName
		}
		if toUpdateUser.Email != "" {
			u.Email = toUpdateUser.Email
		}
	} else {
		u.FirstName = toUpdateUser.FirstName
		u.LastName = toUpdateUser.LastName
		u.Email = toUpdateUser.Email
	}

	return u, nil
}

func (u UserMock) Delete(db database.DatabaseClient) *rest_errors.RestErr {
	if !u.CanDelete {
		return rest_errors.NewInternalServerError("Failed to delete user")
	}

	return nil
}

func (u UserMock) Search(db database.DatabaseClient) ([]users.UserInterface, *rest_errors.RestErr) {
	// TODO
	return nil, nil
}
