package database

import (
	"context"
	"errors"
	"reflect"

	"github.com/ericbg27/top10movies-api/src/datasources/database"
)

type DatabaseClientMock struct {
	Connected      bool
	CanQuery       bool
	CanQueryRow    bool
	CanExec        bool
	CanScanResults bool
}

type ModificationResultMock struct {
	affectedRows int64
}

type UsersSingleElementResultMock struct {
	result  interface{}
	CanScan bool
}

type UsersMultipleElementsResultMock struct {
	results   []interface{}
	scanIndex int
	CanScan   bool
}

func (d *DatabaseClientMock) SetupDbConnection() {
	d.Connected = true
}

func (d *DatabaseClientMock) CloseDbConnection(ctx context.Context) {
	d.Connected = false
}

func (d *DatabaseClientMock) Query(ctx context.Context, query string, arguments ...interface{}) (database.MultipleElementsResult, error) {
	if !d.CanQuery {
		return nil, errors.New("unable to query")
	}

	var usersResult []interface{}

	usersResult = append(usersResult, struct {
		ID          int64
		FirstName   string
		LastName    string
		Email       string
		DateCreated string
		Status      string
		Password    string
	}{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "johndoe@gmail.com",
		DateCreated: "",
		Status:      "",
		Password:    "1234",
	})
	usersResult = append(usersResult, struct {
		ID        int64
		FirstName string
		LastName  string
		Email     string
		Status    string
		Password  string
	}{
		ID:        2,
		FirstName: "Josh",
		LastName:  "Davis",
		Email:     "joshdavis@gmail.com",
		Status:    "active",
		Password:  "12345",
	})

	var result UsersMultipleElementsResultMock
	result.results = usersResult
	result.scanIndex = 0
	result.CanScan = d.CanScanResults

	return &result, nil
}

func (d *DatabaseClientMock) QueryRow(ctx context.Context, query string, arguments ...interface{}) (database.SingleElementResult, error) {
	if !d.CanQueryRow {
		return nil, errors.New("unable to query row")
	}

	userResult := struct {
		ID        int64
		FirstName string
		LastName  string
		Email     string
		Status    string
		Password  string
	}{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@gmail.com",
		Status:    "active",
		Password:  "1234",
	}

	var result UsersSingleElementResultMock
	result.result = userResult
	result.CanScan = d.CanScanResults

	return result, nil
}

func (d *DatabaseClientMock) Exec(ctx context.Context, query string, arguments ...interface{}) (database.ModificationResult, error) {
	if !d.CanExec {
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
	if !us.CanScan {
		return errors.New("failed to scan")
	}

	resultReflection := reflect.TypeOf(us.result)
	resultReflectionValue := reflect.ValueOf(us.result)

	if resultReflection.NumField() < len(arguments) {
		return errors.New("wrong number of arguments provided in scan")
	}

	for index, arg := range arguments {
		if reflect.TypeOf(arg).Elem() != resultReflection.Field(index).Type {
			return errors.New("could not assert type of argument in scan")
		}

		newValue := reflect.ValueOf(resultReflectionValue.FieldByName(resultReflection.Field(index).Name).Interface())
		argCurrentValue := reflect.ValueOf(arguments[index]).Elem()

		argCurrentValue.Set(newValue)
	}

	return nil
}

func (um *UsersMultipleElementsResultMock) Scan(arguments ...interface{}) error {
	if !um.CanScan {
		return errors.New("failed to scan")
	}

	if !um.Next() {
		return errors.New("cannot scan result anymore")
	}

	resultToScan := um.results[um.scanIndex]

	resultReflection := reflect.TypeOf(resultToScan)
	resultReflectionValue := reflect.ValueOf(resultToScan)

	if resultReflection.NumField() < len(arguments) {
		return errors.New("wrong number of arguments provided in scan")
	}

	for index, arg := range arguments {
		if reflect.TypeOf(arg).Elem() != resultReflection.Field(index).Type {
			return errors.New("could not assert type of argument in scan")
		}

		newValue := reflect.ValueOf(resultReflectionValue.FieldByName(resultReflection.Field(index).Name).Interface())
		argCurrentValue := reflect.ValueOf(arguments[index]).Elem()

		argCurrentValue.Set(newValue)
	}

	um.scanIndex++

	return nil
}

func (um *UsersMultipleElementsResultMock) Next() bool {
	return um.scanIndex < len(um.results)
}
