package movies

import "github.com/ericbg27/top10movies-api/src/utils/rest_errors"

type MovieInterface interface {
	AddMovie() *rest_errors.RestErr
	GetMovie() (MovieInterface, *rest_errors.RestErr)
}

type Movie struct {
	ID            int64  `json:"id"`
	OriginalTitle string `json:"original_title"`
	Adult         bool   `json:"adult"`
	ReleaseDate   string `json:"release_date"`
	CreatedAt     string `json:"created_at"`
	Title         string `json:"title"`
	Overview      string `json:"overview"`
}
