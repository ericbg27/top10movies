package database

import "context"

type ModificationResult interface {
	RowsAffected() int64
}

type SingleElementResult interface {
	Scan(...interface{}) error
}

type MultipleElementsResult interface {
	Next() bool
	Scan(...interface{}) error
}

type DatabaseClient interface {
	SetupDbConnection()
	CloseDbConnection(ctx context.Context)
	Query(ctx context.Context, query string, arguments ...interface{}) (MultipleElementsResult, error)
	QueryRow(ctx context.Context, query string, arguments ...interface{}) (SingleElementResult, error)
	Exec(ctx context.Context, query string, arguments ...interface{}) (ModificationResult, error)
}
