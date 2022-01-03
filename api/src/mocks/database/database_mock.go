package database

import (
	"context"
	"errors"
	"reflect"

	"github.com/ericbg27/top10movies-api/src/datasources/database"
	"github.com/ericbg27/top10movies-api/src/domain/users"
)

type DatabaseClientMock struct {
	Connected bool
}

type ModificationResultMock struct {
	affectedRows int64
}

type UsersSingleElementResultMock struct {
	result users.User
}

type UsersMultipleElementsResultMock struct {
	results   []users.User
	scanIndex int
}

func (d *DatabaseClientMock) SetupDbConnection() {
	d.Connected = true
}

func (d *DatabaseClientMock) CloseDbConnection(ctx context.Context) {
	d.Connected = false
}

func (d *DatabaseClientMock) Query(ctx context.Context, query string, arguments interface{}) (database.MultipleElementsResult, error) {
	if query == "error" {
		return nil, errors.New("unable to query")
	}

	var usersResult []users.User

	usersResult = append(usersResult, users.User{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "johndoe@gmail.com",
		DateCreated: "",
		Status:      "",
		Password:    "1234",
	})
	usersResult = append(usersResult, users.User{
		ID:          2,
		FirstName:   "Josh",
		LastName:    "Davis",
		Email:       "joshdavisgmail.com",
		DateCreated: "",
		Status:      "",
		Password:    "12345",
	})

	var result UsersMultipleElementsResultMock
	result.results = usersResult
	result.scanIndex = 0

	return &result, nil
}

func (d *DatabaseClientMock) QueryRow(ctx context.Context, query string, arguments interface{}) (database.SingleElementResult, error) {
	if query == "error" {
		return nil, errors.New("unable to query row")
	}

	userResult := users.User{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "johndoe@gmail.com",
		DateCreated: "",
		Status:      "",
		Password:    "1234",
	}

	var result UsersSingleElementResultMock
	result.result = userResult

	return result, nil
}

func (d *DatabaseClientMock) Exec(ctx context.Context, query string, arguments interface{}) (database.ModificationResult, error) {
	if query == "error" {
		return nil, errors.New("unable to exec")
	}

	var result ModificationResultMock
	result.affectedRows = 1

	return result, nil
}

func (m ModificationResultMock) RowsAffected() int64 {
	return m.affectedRows
}

func (us UsersSingleElementResultMock) Scan(arguments ...interface{}) error {
	resultReflection := reflect.TypeOf(us.result)
	resultReflectionValue := reflect.ValueOf(us.result)

	if resultReflection.NumField() != len(arguments) {
		return errors.New("wrong number of arguments provided in scan")
	}

	for index, arg := range arguments {
		if reflect.TypeOf(arg) != resultReflection.Field(index).Type {
			return errors.New("could not assert type of argument in scan")
		}

		arguments[index] = resultReflectionValue.FieldByName(resultReflection.Field(index).Name).Interface()
	}

	return nil
}

func (um *UsersMultipleElementsResultMock) Scan(arguments ...interface{}) error {
	if !um.Next() {
		return errors.New("cannot scan result anymore")
	}

	resultToScan := um.results[um.scanIndex]

	resultReflection := reflect.TypeOf(resultToScan)
	resultReflectionValue := reflect.ValueOf(resultToScan)

	if resultReflection.NumField() != len(arguments) {
		return errors.New("wrong number of arguments provided in scan")
	}

	for index, arg := range arguments {
		if reflect.TypeOf(arg) != resultReflection.Field(index).Type {
			return errors.New("could not assert type of argument in scan")
		}

		arguments[index] = resultReflectionValue.FieldByName(resultReflection.Field(index).Name).Interface()
	}

	um.scanIndex++

	return nil
}

func (um *UsersMultipleElementsResultMock) Next() bool {
	return um.scanIndex < len(um.results)
}
