package authorization

import (
	"errors"
	"strconv"
	"strings"

	auth "github.com/ericbg27/top10movies-api/src/utils/authorization"
)

type AuthorizationMock struct {
	CanCreate bool
}

func (a AuthorizationMock) CreateToken(userId int64) (*auth.TokenDetails, error) {
	if !a.CanCreate {
		return nil, errors.New("failed to create token")
	}

	var sb strings.Builder

	sb.WriteString("token_")
	sb.WriteString(strconv.Itoa(int(userId)))

	s := sb.String()

	tokenInfo := &auth.TokenDetails{
		AccessToken: s,
	}

	return tokenInfo, nil
}

func (a AuthorizationMock) FetchAuth(bearToken string) (uint64, error) {
	// TODO
	return 0, nil
}
