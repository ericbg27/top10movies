package movies

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/ericbg27/top10movies-api/src/utils/logger"
	"github.com/ericbg27/top10movies-api/src/utils/rest_errors"
)

const (
	CreatedAtLayout = "2006-02-01T15:04:05Z"
)

var (
	cachettl = config.GetConfig().Redis.CacheTtl
)

func (m MovieInfo) AddMovie() *rest_errors.RestErr {
	m.CreatedAt = time.Now().Format(CreatedAtLayout)

	marshelledMovie, err := json.Marshal(m)
	if err != nil {
		logger.Error("Error when trying to add movie", err)
		return rest_errors.NewInternalServerError("Error when trying to add movie")
	}

	var movieRedisKey strings.Builder
	movieRedisKey.WriteString("movie:")
	movieRedisKey.WriteString(strconv.Itoa(m.Movie.ID))
	redisdb.Client.Set(movieRedisKey.String(), marshelledMovie, time.Duration(cachettl*int64(time.Minute)))

	return nil
}

func (m MovieInfo) GetMovie() (MovieInterface, *rest_errors.RestErr) {
	var savedMovie MovieInfo

	var movieRedisKey strings.Builder
	movieRedisKey.WriteString("movie:")
	movieRedisKey.WriteString(strconv.Itoa(m.Movie.ID))

	result, err := redisdb.Client.Get(movieRedisKey.String()).Result()
	if err != nil && err != redisdb.RedisNil {
		logger.Error("Error when trying to get movie", err)
		return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
	} else if err == redisdb.RedisNil {
		savedMovie.Movie.ID = -1
	} else {
		marshalErr := json.Unmarshal([]byte(result), &savedMovie)
		if marshalErr != nil {
			logger.Error("Error when trying to get movie", marshalErr)
			return nil, rest_errors.NewInternalServerError("Error when trying to get movie")
		}
	}

	return savedMovie, nil
}
