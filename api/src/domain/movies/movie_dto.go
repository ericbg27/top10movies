package movies

import (
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
	"github.com/ryanbradynd05/go-tmdb"

	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
)

type MovieInterface interface {
	AddMovie(redisClient *redisdb.RedisClient) *rest_errors.RestErr
	GetMovie(redisClient *redisdb.RedisClient) (MovieInterface, *rest_errors.RestErr)
}

type MovieInfo struct {
	Movie     tmdb.Movie `json:"movie_info"`
	CreatedAt string     `json:"created_at"`
}
