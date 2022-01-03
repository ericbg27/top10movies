package postgresdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/ericbg27/top10movies-api/src/datasources/database"
	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDBClient struct {
	Client *pgxpool.Conn
}

var (
	host     = config.GetConfig().Database.Host
	port     = config.GetConfig().Database.Port
	user     = config.GetConfig().Database.User
	password = config.GetConfig().Database.Password
	dbname   = config.GetConfig().Database.DbName
	loglevel = config.GetConfig().Database.LogLevel
)

func (p *PostgresDBClient) SetupDbConnection() {
	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname))
	if err != nil {
		logger.Error("Error when parsing database connection string", err)
		panic(err)
	}

	config.ConnConfig.Logger = logger.GetLogger()

	level, err := pgx.LogLevelFromString(strings.ToLower(loglevel))
	if err != nil {
		level = pgx.LogLevelInfo
	}
	config.ConnConfig.LogLevel = level

	clientPool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to database: %v\n", err.Error()), err)
		panic(err)
	}

	p.Client, err = clientPool.Acquire(context.Background())
	if err != nil {
		logger.Error("Unable to acquire connection to database", err)
		panic(err)
	}
}

func (p *PostgresDBClient) CloseDbConnection(ctx context.Context) {
	p.Client.Conn().Close(context.Background())
}

func (p *PostgresDBClient) Query(ctx context.Context, query string, arguments ...interface{}) (database.MultipleElementsResult, error) {
	result, err := p.Client.Query(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostgresDBClient) QueryRow(ctx context.Context, query string, arguments ...interface{}) (database.SingleElementResult, error) {
	result := p.Client.QueryRow(ctx, query, arguments...)

	return result, nil
}

func (p *PostgresDBClient) Exec(ctx context.Context, query string, arguments ...interface{}) (database.ModificationResult, error) {
	result, err := p.Client.Exec(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
